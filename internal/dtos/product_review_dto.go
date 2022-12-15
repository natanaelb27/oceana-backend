package dtos

import (
	"time"

	"github.com/google/uuid"
)

type ProductReviewList struct {
	ID                uint
	UserID            uuid.UUID
	FullName          string
	Username          string
	Photo             *string
	ProductID         uint
	Rating            uint
	Review            string
	ProductVariations string
	CreatedAt         time.Time
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
