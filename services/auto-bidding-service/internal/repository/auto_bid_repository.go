package repository

import (
	"auto-bidding-service/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
)

// AutoBidRepository xử lý thao tác với database cho auto-bids
type AutoBidRepository struct {
	db *pg.DB
}

// NewAutoBidRepository tạo repository mới
func NewAutoBidRepository(db *pg.DB) *AutoBidRepository {
	return &AutoBidRepository{db: db}
}

// Create tạo một auto-bid mới
func (r *AutoBidRepository) Create(ctx context.Context, autoBid *models.AutoBid) error {
	autoBid.CreatedAt = time.Now()
	autoBid.UpdatedAt = time.Now()
	autoBid.Status = models.AutoBidStatusActive

	_, err := r.db.ModelContext(ctx, autoBid).Insert()
	if err != nil {
		return fmt.Errorf("failed to create auto-bid: %w", err)
	}
	return nil
}

// Update cập nhật auto-bid
func (r *AutoBidRepository) Update(ctx context.Context, autoBid *models.AutoBid) error {
	autoBid.UpdatedAt = time.Now()

	_, err := r.db.ModelContext(ctx, autoBid).
		WherePK().
		Update()
	if err != nil {
		return fmt.Errorf("failed to update auto-bid: %w", err)
	}
	return nil
}

// GetByID lấy auto-bid theo ID
func (r *AutoBidRepository) GetByID(ctx context.Context, id int64) (*models.AutoBid, error) {
	autoBid := &models.AutoBid{ID: id}
	err := r.db.ModelContext(ctx, autoBid).WherePK().Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("auto-bid not found")
		}
		return nil, fmt.Errorf("failed to get auto-bid: %w", err)
	}
	return autoBid, nil
}

// GetActiveByProduct lấy tất cả auto-bid ACTIVE của một sản phẩm, sắp xếp theo max_amount giảm dần
func (r *AutoBidRepository) GetActiveByProduct(ctx context.Context, productID int64) ([]*models.AutoBid, error) {
	var autoBids []*models.AutoBid
	err := r.db.ModelContext(ctx, &autoBids).
		Where("product_id = ?", productID).
		Where("status = ?", models.AutoBidStatusActive).
		Order("max_amount DESC", "created_at ASC"). // Sắp xếp theo max_amount giảm, nếu bằng nhau thì người tạo trước win
		Select()

	if err != nil {
		return nil, fmt.Errorf("failed to get active auto-bids: %w", err)
	}
	return autoBids, nil
}

// GetByBidder lấy tất cả auto-bid của một bidder
func (r *AutoBidRepository) GetByBidder(ctx context.Context, bidderID int64) ([]*models.AutoBid, error) {
	var autoBids []*models.AutoBid
	err := r.db.ModelContext(ctx, &autoBids).
		Where("bidder_id = ?", bidderID).
		Order("created_at DESC").
		Select()

	if err != nil {
		return nil, fmt.Errorf("failed to get bidder's auto-bids: %w", err)
	}
	return autoBids, nil
}

// GetByBidderAndProduct lấy auto-bid của một bidder cho một sản phẩm cụ thể
func (r *AutoBidRepository) GetByBidderAndProduct(ctx context.Context, bidderID, productID int64) (*models.AutoBid, error) {
	autoBid := &models.AutoBid{}
	err := r.db.ModelContext(ctx, autoBid).
		Where("bidder_id = ?", bidderID).
		Where("product_id = ?", productID).
		Where("status = ?", models.AutoBidStatusActive).
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil // Không có auto-bid active
		}
		return nil, fmt.Errorf("failed to get auto-bid: %w", err)
	}
	return autoBid, nil
}

// UpdateStatus cập nhật trạng thái của auto-bid
func (r *AutoBidRepository) UpdateStatus(ctx context.Context, id int64, status models.AutoBidStatus) error {
	_, err := r.db.ModelContext(ctx, (*models.AutoBid)(nil)).
		Set("status = ?", status).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Update()

	if err != nil {
		return fmt.Errorf("failed to update auto-bid status: %w", err)
	}
	return nil
}

// UpdateCurrentAmount cập nhật current_amount của auto-bid
func (r *AutoBidRepository) UpdateCurrentAmount(ctx context.Context, id int64, amount float64) error {
	_, err := r.db.ModelContext(ctx, (*models.AutoBid)(nil)).
		Set("current_amount = ?", amount).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Update()

	if err != nil {
		return fmt.Errorf("failed to update current amount: %w", err)
	}
	return nil
}

// DeactivateOldAutoBids đánh dấu các auto-bid cũ của bidder cho sản phẩm là CANCELLED
func (r *AutoBidRepository) DeactivateOldAutoBids(ctx context.Context, bidderID, productID int64) error {
	_, err := r.db.ModelContext(ctx, (*models.AutoBid)(nil)).
		Set("status = ?", models.AutoBidStatusCancelled).
		Set("updated_at = ?", time.Now()).
		Where("bidder_id = ?", bidderID).
		Where("product_id = ?", productID).
		Where("status = ?", models.AutoBidStatusActive).
		Update()

	if err != nil {
		return fmt.Errorf("failed to deactivate old auto-bids: %w", err)
	}
	return nil
}
