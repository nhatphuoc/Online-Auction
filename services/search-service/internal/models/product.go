package models

import "time"

// Product represents a product in PostgreSQL
type Product struct {
	ID            int64      `json:"id" pg:"id,pk"`
	Name          string     `json:"name" pg:"name"`
	Description   string     `json:"description" pg:"description"`
	CategoryID    int64      `json:"category_id" pg:"category_id"`
	SellerID      int64      `json:"seller_id" pg:"seller_id"`
	StartingPrice float64    `json:"starting_price" pg:"starting_price"`
	CurrentPrice  float64    `json:"current_price" pg:"current_price"`
	BuyNowPrice   *float64   `json:"buy_now_price,omitempty" pg:"buy_now_price"`
	StepPrice     float64    `json:"step_price" pg:"step_price"`
	Status        string     `json:"status" pg:"status"`
	ThumbnailURL  string     `json:"thumbnail_url" pg:"thumbnail_url"`
	AutoExtend    bool       `json:"auto_extend" pg:"auto_extend"`
	CurrentBidder *int64     `json:"current_bidder,omitempty" pg:"current_bidder"`
	EndAt         time.Time  `json:"end_at" pg:"end_at"`
	CreatedAt     time.Time  `json:"created_at" pg:"created_at"`
}

// ProductESDocument represents product document in Elasticsearch
type ProductESDocument struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	NameNoAccent     string    `json:"name_no_accent"`
	Description      string    `json:"description"`
	DescriptionNoAccent string `json:"description_no_accent"`
	CategoryID       int64     `json:"category_id"`
	CategoryName     string    `json:"category_name"`
	CategorySlug     string    `json:"category_slug"`
	SellerID         int64     `json:"seller_id"`
	StartingPrice    float64   `json:"starting_price"`
	CurrentPrice     float64   `json:"current_price"`
	BuyNowPrice      *float64  `json:"buy_now_price,omitempty"`
	StepPrice        float64   `json:"step_price"`
	Status           string    `json:"status"`
	ThumbnailURL     string    `json:"thumbnail_url"`
	AutoExtend       bool      `json:"auto_extend"`
	CurrentBidder    *int64    `json:"current_bidder,omitempty"`
	EndAt            time.Time `json:"end_at"`
	CreatedAt        time.Time `json:"created_at"`
}

// SearchRequest represents search parameters
type SearchRequest struct {
	Query      string   `json:"query" query:"query"`
	CategoryID *int64   `json:"category_id" query:"category_id"`
	Status     string   `json:"status" query:"status"`
	MinPrice   *float64 `json:"min_price" query:"min_price"`
	MaxPrice   *float64 `json:"max_price" query:"max_price"`
	SortBy     string   `json:"sort_by" query:"sort_by"`
	SortOrder  string   `json:"sort_order" query:"sort_order"`
	Page       int      `json:"page" query:"page"`
	PageSize   int      `json:"page_size" query:"page_size"`
}

// SearchResponse represents search results
type SearchResponse struct {
	Products   []ProductESDocument `json:"products"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}
