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
