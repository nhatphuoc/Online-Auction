package models

import "time"

// Comment represents a comment on a product
type Comment struct {
	tableName struct{} `pg:"comments"`

	ID        int       `json:"id" pg:"id,pk"`
	ProductID int       `json:"product_id" pg:"product_id,notnull"`
	SenderID  int       `json:"sender_id" pg:"sender_id,notnull"`
	Content   string    `json:"content" pg:"content,notnull"`
	CreatedAt time.Time `json:"created_at" pg:"created_at"`
}

// CommentRequest represents the request payload for sending a comment
type CommentRequest struct {
	ProductID int    `json:"product_id" validate:"required"`
	Content   string `json:"content" validate:"required"`
}

// CommentResponse includes sender information
type CommentResponse struct {
	ID         int       `json:"id"`
	ProductID  int       `json:"product_id"`
	SenderID   int       `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// WebSocketMessage represents message sent through WebSocket
type WebSocketMessage struct {
	Type      string      `json:"type"` // join, leave, comment, typing
	ProductID int         `json:"product_id"`
	Content   string      `json:"content,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}
