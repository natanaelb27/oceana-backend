package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductReviewMedia struct {
	ID              uint
	ProductReviewID uint
	ProductReview   ProductReview
	MediaURL        string
	Type            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}
