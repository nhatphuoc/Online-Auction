package service

import (
	"auto-bidding-service/internal/client"
	"auto-bidding-service/internal/models"
	"auto-bidding-service/internal/repository"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// AutoBidService xử lý logic nghiệp vụ cho auto-bidding
type AutoBidService struct {
	repo                 *repository.AutoBidRepository
	biddingServiceClient *client.BiddingServiceClient
	productServiceClient *client.ProductServiceClient
}

// NewAutoBidService tạo service mới
func NewAutoBidService(
	repo *repository.AutoBidRepository,
	biddingServiceClient *client.BiddingServiceClient,
	productServiceClient *client.ProductServiceClient,
) *AutoBidService {
	return &AutoBidService{
		repo:                 repo,
		biddingServiceClient: biddingServiceClient,
		productServiceClient: productServiceClient,
	}
}

// CreateAutoBid tạo một auto-bid mới cho bidder
func (s *AutoBidService) CreateAutoBid(ctx context.Context, bidderID, productID int64, maxAmount float64, userToken string) (*models.AutoBid, error) {
	// 1. Kiểm tra sản phẩm có tồn tại và đang active không
	product, err := s.productServiceClient.GetProduct(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product info: %w", err)
	}

	if product.Status != "ACTIVE" {
		return nil, fmt.Errorf("product is not active for bidding")
	}

	// 2. Kiểm tra maxAmount phải lớn hơn giá hiện tại
	if maxAmount <= product.CurrentPrice {
		return nil, fmt.Errorf("max amount must be greater than current price %.2f", product.CurrentPrice)
	}

	// 3. Deactivate auto-bid cũ của bidder cho sản phẩm này (nếu có)
	if err := s.repo.DeactivateOldAutoBids(ctx, bidderID, productID); err != nil {
		slog.Error("Failed to deactivate old auto-bids", "error", err)
		// Không return error, tiếp tục tạo mới
	}

	// 4. Tạo auto-bid mới
	autoBid := &models.AutoBid{
		ProductID:     productID,
		BidderID:      bidderID,
		MaxAmount:     maxAmount,
		CurrentAmount: 0, // Chưa bid
		Status:        models.AutoBidStatusActive,
	}

	if err := s.repo.Create(ctx, autoBid); err != nil {
		return nil, fmt.Errorf("failed to create auto-bid: %w", err)
	}

	// 5. Trigger auto-bidding ngay lập tức
	go s.TriggerAutoBidding(context.Background(), productID, product.CurrentPrice, product.StepPrice, bidderID, maxAmount, userToken)

	return autoBid, nil
}

// TriggerAutoBidding xử lý logic auto-bidding khi có bid mới
// Logic:
// - Lấy tất cả auto-bid ACTIVE của sản phẩm, sắp xếp theo max_amount giảm dần
// - Những người có max_amount < giá hiện tại → OUTBID
// - Những người có max_amount >= giá hiện tại:
//   - Người thứ 2,3,4... bid hết max của họ
//   - Người thứ nhất (max cao nhất) chỉ bid cao hơn người thứ 2 một bước giá
func (s *AutoBidService) TriggerAutoBidding(ctx context.Context, productID int64, currentPrice, stepPrice float64, triggerBidderID int64, triggerAmount float64, userToken string) error {
	slog.Info("Triggering auto-bidding",
		"product_id", productID,
		"current_price", currentPrice,
		"step_price", stepPrice,
		"trigger_bidder", triggerBidderID,
		"trigger_amount", triggerAmount)

	// 1. Lấy tất cả auto-bid ACTIVE, sắp xếp theo max_amount DESC
	autoBids, err := s.repo.GetActiveByProduct(ctx, productID)
	if err != nil {
		slog.Error("Failed to get active auto-bids", "error", err)
		return err
	}

	if len(autoBids) == 0 {
		slog.Info("No active auto-bids found for product", "product_id", productID)
		return nil
	}

	slog.Info("Found active auto-bids", "count", len(autoBids))

	// 2. Lọc ra các auto-bid có max_amount > currentPrice (có khả năng bid)
	var eligibleBids []*models.AutoBid
	for _, ab := range autoBids {
		if ab.MaxAmount > currentPrice {
			eligibleBids = append(eligibleBids, ab)
		} else {
			// Đánh dấu OUTBID cho những người có max < giá hiện tại
			s.repo.UpdateStatus(ctx, ab.ID, models.AutoBidStatusOutbid)
			slog.Info("Auto-bid marked as OUTBID", "auto_bid_id", ab.ID, "max_amount", ab.MaxAmount)
		}
	}

	if len(eligibleBids) == 0 {
		slog.Info("No eligible auto-bids after filtering")
		return nil
	}

	// 3. Xử lý logic bidding
	// - Nếu chỉ có 1 người: bid = currentPrice + stepPrice
	// - Nếu có 2+ người:
	//   + Người 2,3,4,... bid hết max của họ
	//   + Người 1 (max cao nhất) bid = min(max của người 2 + stepPrice, max của người 1)

	if len(eligibleBids) == 1 {
		// Chỉ có 1 người, bid ngay
		autoBid := eligibleBids[0]
		nextBidAmount := currentPrice + stepPrice

		// Không vượt quá max
		if nextBidAmount > autoBid.MaxAmount {
			nextBidAmount = autoBid.MaxAmount
		}

		s.executeBid(ctx, autoBid, nextBidAmount, userToken)
	} else {
		// Có nhiều người
		highestAutoBid := eligibleBids[0] // Người có max cao nhất
		secondHighest := eligibleBids[1]  // Người thứ 2

		// Người từ thứ 2 trở đi: bid hết max của họ
		for i := len(eligibleBids) - 1; i >= 1; i-- {
			autoBid := eligibleBids[i]
			s.executeBid(ctx, autoBid, autoBid.MaxAmount, userToken)

			// Delay một chút để tránh race condition
			time.Sleep(100 * time.Millisecond)
		}

		// Người thứ nhất: bid cao hơn người thứ 2 một bước giá
		winningBidAmount := secondHighest.MaxAmount + stepPrice

		// Không vượt quá max của người thứ nhất
		if winningBidAmount > highestAutoBid.MaxAmount {
			winningBidAmount = highestAutoBid.MaxAmount
		}

		// Delay trước khi bid cuối cùng
		time.Sleep(100 * time.Millisecond)
		s.executeBid(ctx, highestAutoBid, winningBidAmount, userToken)
	}

	return nil
}

// executeBid thực hiện việc đặt giá qua bidding-service
func (s *AutoBidService) executeBid(ctx context.Context, autoBid *models.AutoBid, amount float64, userToken string) {
	requestID := uuid.New().String()

	slog.Info("Executing auto-bid",
		"auto_bid_id", autoBid.ID,
		"bidder_id", autoBid.BidderID,
		"product_id", autoBid.ProductID,
		"amount", amount,
		"max_amount", autoBid.MaxAmount,
		"request_id", requestID)

	// Gọi bidding-service để đặt giá
	resp, err := s.biddingServiceClient.PlaceBid(
		autoBid.ProductID,
		autoBid.BidderID,
		amount,
		requestID,
		userToken,
	)

	if err != nil {
		slog.Error("Failed to place bid via bidding-service",
			"error", err,
			"auto_bid_id", autoBid.ID)
		return
	}

	if resp.Success {
		// Cập nhật current_amount
		s.repo.UpdateCurrentAmount(ctx, autoBid.ID, amount)
		slog.Info("Auto-bid executed successfully",
			"auto_bid_id", autoBid.ID,
			"amount", amount)
	} else {
		slog.Error("Bid rejected by bidding-service",
			"auto_bid_id", autoBid.ID,
			"message", resp.Message)

		// Nếu bid thất bại, có thể cần đánh dấu auto-bid
		// Tùy theo lý do thất bại mà xử lý khác nhau
	}
}

// GetAutoBidsByBidder lấy danh sách auto-bid của một bidder
func (s *AutoBidService) GetAutoBidsByBidder(ctx context.Context, bidderID int64) ([]*models.AutoBid, error) {
	return s.repo.GetByBidder(ctx, bidderID)
}

// GetAutoBidByID lấy thông tin auto-bid theo ID
func (s *AutoBidService) GetAutoBidByID(ctx context.Context, id int64) (*models.AutoBid, error) {
	return s.repo.GetByID(ctx, id)
}

// CancelAutoBid hủy một auto-bid
func (s *AutoBidService) CancelAutoBid(ctx context.Context, id, bidderID int64) error {
	// Kiểm tra auto-bid có thuộc về bidder không
	autoBid, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if autoBid.BidderID != bidderID {
		return fmt.Errorf("unauthorized: auto-bid does not belong to bidder")
	}

	if autoBid.Status != models.AutoBidStatusActive {
		return fmt.Errorf("auto-bid is not active, cannot cancel")
	}

	return s.repo.UpdateStatus(ctx, id, models.AutoBidStatusCancelled)
}
