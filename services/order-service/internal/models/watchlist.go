package models

import "time"

// WatchList represents a user's favorite/watched product
type WatchList struct {
	tableName struct{} `pg:"watch_list"`

	ID        int64     `json:"id" pg:"id,pk"`
	UserID    int64     `json:"user_id" pg:"user_id,notnull"`       // ID người dùng
	ProductID int64     `json:"product_id" pg:"product_id,notnull"` // ID sản phẩm đấu giá
	CreatedAt time.Time `json:"created_at" pg:"created_at,default:now()"`
}

// AddToWatchListRequest represents request to add product to watch list
type AddToWatchListRequest struct {
	ProductID int64 `json:"product_id" validate:"required"`
}

// WatchListResponse represents watch list item with product details
type WatchListResponse struct {
	ID           int64     `json:"id" pg:"id"`
	ProductID    int64     `json:"product_id" pg:"product_id"`
	ThumbnailURL string    `json:"thumbnailUrl" pg:"thumbnail_url"`
	Name         string    `json:"name" pg:"name"`
	CurrentPrice float64   `json:"currentPrice" pg:"current_price"`
	BuyNowPrice  float64   `json:"buyNowPrice" pg:"buy_now_price"`
	CreatedAt    time.Time `json:"createdAt" pg:"created_at"`
	EndAt        time.Time `json:"endAt" pg:"end_at"`
	BidCount     int64     `json:"bidCount" pg:"bid_count"`
	CategoryName string    `json:"categoryName" pg:"category_name"`
}
