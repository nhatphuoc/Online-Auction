package models

import "time"

// AutoBidStatus đại diện cho trạng thái của auto-bid
type AutoBidStatus string

const (
	AutoBidStatusActive    AutoBidStatus = "ACTIVE"     // Đang hoạt động
	AutoBidStatusWon       AutoBidStatus = "WON"        // Đã thắng
	AutoBidStatusOutbid    AutoBidStatus = "OUTBID"     // Bị đấu giá vượt mức (max_amount < giá hiện tại)
	AutoBidStatusCancelled AutoBidStatus = "CANCELLED"  // Đã hủy
	AutoBidStatusExpired   AutoBidStatus = "EXPIRED"    // Hết hạn (sản phẩm kết thúc đấu giá)
)

// AutoBid đại diện cho một lệnh đấu giá tự động
type AutoBid struct {
	ID              int64         `db:"id" json:"id"`
	ProductID       int64         `db:"product_id" json:"product_id"`
	BidderID        int64         `db:"bidder_id" json:"bidder_id"`
	MaxAmount       float64       `db:"max_amount" json:"max_amount"`           // Giá tối đa mà bidder sẵn sàng trả
	CurrentAmount   float64       `db:"current_amount" json:"current_amount"`   // Giá hiện tại đã bid
	Status          AutoBidStatus `db:"status" json:"status"`
	CreatedAt       time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at" json:"updated_at"`
}

// CreateAutoBidRequest là request để tạo auto-bid mới
type CreateAutoBidRequest struct {
	ProductID int64   `json:"product_id" validate:"required,gt=0"`
	MaxAmount float64 `json:"max_amount" validate:"required,gt=0"`
}

// AutoBidResponse là response cho auto-bid
type AutoBidResponse struct {
	ID            int64         `json:"id"`
	ProductID     int64         `json:"product_id"`
	BidderID      int64         `json:"bidder_id"`
	MaxAmount     float64       `json:"max_amount"`
	CurrentAmount float64       `json:"current_amount"`
	Status        AutoBidStatus `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// TriggerAutoBidRequest là request để trigger auto-bidding khi có bid mới
type TriggerAutoBidRequest struct {
	ProductID     int64   `json:"product_id" validate:"required,gt=0"`
	CurrentPrice  float64 `json:"current_price" validate:"required,gt=0"`
	BidIncrement  float64 `json:"bid_increment" validate:"required,gt=0"`
	NewBidderID   int64   `json:"new_bidder_id" validate:"required,gt=0"` // ID của người vừa bid
	NewBidAmount  float64 `json:"new_bid_amount" validate:"required,gt=0"` // Số tiền của bid mới
}

// ProductInfo lưu thông tin sản phẩm cần thiết cho auto-bidding
type ProductInfo struct {
	ID            int64   `json:"id"`
	CurrentPrice  float64 `json:"current_price"`
	BidIncrement  float64 `json:"bid_increment"`
	HighestBidder int64   `json:"highest_bidder"`
}
