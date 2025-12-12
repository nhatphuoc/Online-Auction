package models

import "time"

// Category represents a category in PostgreSQL
type Category struct {
	ID           int64     `json:"id" pg:"id,pk"`
	Name         string    `json:"name" pg:"name"`
	Slug         string    `json:"slug" pg:"slug"`
	Description  string    `json:"description" pg:"description"`
	ParentID     *int64    `json:"parent_id,omitempty" pg:"parent_id"`
	Level        int       `json:"level" pg:"level"`
	IsActive     bool      `json:"is_active" pg:"is_active"`
	DisplayOrder int       `json:"display_order" pg:"display_order"`
	CreatedAt    time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" pg:"updated_at"`
}

// CategoryESDocument represents category document in Elasticsearch
type CategoryESDocument struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	NameNoAccent string    `json:"name_no_accent"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	ParentID     *int64    `json:"parent_id,omitempty"`
	Level        int       `json:"level"`
	IsActive     bool      `json:"is_active"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
