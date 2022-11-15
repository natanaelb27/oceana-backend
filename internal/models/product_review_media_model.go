package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductReviewMedia struct {
	ID              uint
	ProductReviewId uint
	ProductReview   ProductReview
	MediaURL        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}
