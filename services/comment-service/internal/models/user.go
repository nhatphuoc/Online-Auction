package models

import "time"

// User represents user table in database
type User struct {
	tableName struct{} `pg:"users"`

	ID                     int       `json:"id" pg:"id,pk"`
	BirthDay               time.Time `json:"birth_day" pg:"birth_day"`
	Email                  string    `json:"email" pg:"email,unique,notnull"`
	EmailVerified          bool      `json:"email_verified" pg:"email_verified"`
	FullName               string    `json:"full_name" pg:"full_name"`
	IsSellerRequestSent    bool      `json:"is_seller_request_sent" pg:"is_seller_request_sent"`
	Password               string    `json:"-" pg:"password"` // Never expose password
	Role                   string    `json:"role" pg:"role"`
	TotalNumberGoodReviews int       `json:"total_number_good_reviews" pg:"total_number_good_reviews,notnull"`
	TotalNumberReviews     int       `json:"total_number_reviews" pg:"total_number_reviews,notnull"`
}
