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

type ProductReviewGetByProductID struct {
	ID                 uint      `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	FullName           string    `json:"full_name"`
	Username           string    `json:"username"`
	Photo              *string   `json:"photo"`
	ProductID          uint      `json:"product_id"`
	Rating             uint      `json:"rating"`
	Review             string    `json:"review"`
	ProductReviewMedia []*string `json:"product_review_media"`
	ProductVariations  string    `json:"product_variations"`
	CreatedAt          time.Time `json:"created_at"`
}

func NewProductReviewGetByProductID(source *ProductReviewList, productReviewMediaLinks []*string) *ProductReviewGetByProductID {
	return &ProductReviewGetByProductID{
		ID:                 source.ID,
		UserID:             source.UserID,
		FullName:           source.FullName,
		Username:           source.Username,
		Photo:              source.Photo,
		ProductID:          source.ProductID,
		Rating:             source.Rating,
		Review:             source.Review,
		ProductVariations:  source.ProductVariations,
		ProductReviewMedia: productReviewMediaLinks,
	}
}
