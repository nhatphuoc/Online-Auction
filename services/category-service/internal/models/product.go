package models

import "time"

// Product represents a product in the auction system
type Product struct {
	tableName struct{} `pg:"products"`

	ID            int64     `json:"id" pg:"id,pk"`
	Name          string    `json:"name" pg:"name,notnull"`
	Description   string    `json:"description" pg:"description"`
	CategoryID    int64     `json:"category_id" pg:"category_id,notnull"`
	SellerID      int64     `json:"seller_id" pg:"seller_id,notnull"`
	StartingPrice float64   `json:"starting_price" pg:"starting_price,notnull"`
	CurrentPrice  float64   `json:"current_price" pg:"current_price"`
	BuyNowPrice   *float64  `json:"buy_now_price,omitempty" pg:"buy_now_price"`
	StepPrice     float64   `json:"step_price" pg:"step_price,notnull"`
	Status        string    `json:"status" pg:"status,notnull"`
	ThumbnailURL  string    `json:"thumbnail_url" pg:"thumbnail_url"`
	AutoExtend    bool      `json:"auto_extend" pg:"auto_extend,notnull,default:false"`
	EndAt         time.Time `json:"end_at" pg:"end_at,notnull"`
	CreatedAt     time.Time `json:"created_at" pg:"created_at,default:now()"`
	
	// Relations
	Category *Category       `json:"category,omitempty" pg:"rel:has-one"`
	Images   []*ProductImage `json:"images,omitempty" pg:"rel:has-many"`
}

// ProductImage represents product images
type ProductImage struct {
	tableName struct{} `pg:"product_images"`

	ProductID int64  `json:"product_id" pg:"product_id,notnull"`
	ImageURL  string `json:"image_url" pg:"image_url,notnull"`
}

// ProductListResponse represents paginated product list
type ProductListResponse struct {
	Products []*Product `json:"products"`
	Total    int        `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	TotalPages int      `json:"total_pages"`
}

// ProductQueryParams represents query parameters for listing products
type ProductQueryParams struct {
	CategoryID int64  `query:"category_id"`
	Status     string `query:"status"`
	Page       int    `query:"page"`
	PageSize   int    `query:"page_size"`
	SortBy     string `query:"sort_by"`
	SortOrder  string `query:"sort_order"`
}
