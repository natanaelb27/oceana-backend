package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ProductReviewList struct {
	ID                 uint           `json:"id"`
	UserID             uuid.UUID      `json:"user_id"`
	FullName           string         `json:"full_name"`
	Username           string         `json:"username"`
	Media              pq.StringArray `json:"media" gorm:"type:text[]"`
	Photo              *string        `json:"photo"`
	ProductID          uint           `json:"product_id"`
	Rating             uint           `json:"rating"`
	Review             string         `json:"review"`
	ProductReviewMedia []*string      `json:"product_review_media"`
	ProductVariations  string         `json:"product_variations"`
	CreatedAt          time.Time      `json:"created_at"`
}

type ProductReviewDetails struct {
	RatingCount         uint    `json:"rating_count"`
	ReviewCount         uint    `json:"review_count"`
	RatingStar          float64 `json:"rating_star"`
	SatisfiedPercentage float64 `json:"satisfied_percentage"`
	Star5               uint    `json:"star_5" gorm:"column:star_5"`
	Star4               uint    `json:"star_4" gorm:"column:star_4"`
	Star3               uint    `json:"star_3" gorm:"column:star_3"`
	Star2               uint    `json:"star_2" gorm:"column:star_2"`
	Star1               uint    `json:"star_1" gorm:"column:star_1"`
}
