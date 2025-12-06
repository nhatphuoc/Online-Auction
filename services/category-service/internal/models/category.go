package models

import "time"

// Category represents a product category with hierarchical structure
type Category struct {
	tableName struct{} `pg:"categories"`

	ID          int64      `json:"id" pg:"id,pk"`
	Name        string     `json:"name" pg:"name,notnull"`
	Slug        string     `json:"slug" pg:"slug,notnull,unique"`
	Description string     `json:"description" pg:"description"`
	ParentID    *int64     `json:"parent_id,omitempty" pg:"parent_id"`
	Level       int        `json:"level" pg:"level,notnull,default:1"`
	IsActive    bool       `json:"is_active" pg:"is_active,notnull,default:true"`
	DisplayOrder int       `json:"display_order" pg:"display_order,default:0"`
	CreatedAt   time.Time  `json:"created_at" pg:"created_at,default:now()"`
	UpdatedAt   time.Time  `json:"updated_at" pg:"updated_at,default:now()"`
	
	// Relations
	Parent      *Category   `json:"parent,omitempty" pg:"rel:has-one,fk:parent_id"`
	Children    []*Category `json:"children,omitempty" pg:"rel:has-many,join_fk:parent_id"`
}

// CreateCategoryRequest represents the request payload for creating a category
type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Slug        string  `json:"slug" validate:"required,min=2,max=255"`
	Description string  `json:"description"`
	ParentID    *int64  `json:"parent_id"`
	DisplayOrder int    `json:"display_order"`
}

// UpdateCategoryRequest represents the request payload for updating a category
type UpdateCategoryRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=255"`
	Slug        *string `json:"slug" validate:"omitempty,min=2,max=255"`
	Description *string `json:"description"`
	ParentID    *int64  `json:"parent_id"`
	IsActive    *bool   `json:"is_active"`
	DisplayOrder *int   `json:"display_order"`
}

// CategoryResponse represents category with children
type CategoryResponse struct {
	ID          int64               `json:"id"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Description string              `json:"description"`
	ParentID    *int64              `json:"parent_id,omitempty"`
	Level       int                 `json:"level"`
	IsActive    bool                `json:"is_active"`
	DisplayOrder int                `json:"display_order"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Children    []*CategoryResponse `json:"children,omitempty"`
}

// CategoryTreeResponse represents the hierarchical category tree
type CategoryTreeResponse struct {
	Categories []*CategoryResponse `json:"categories"`
}
