package dtos

import (
	"mime/multipart"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type ProductReviewCreateRequest struct {
	UserID        uuid.UUID
	OrderDetailID uuid.UUID
	ProductID     uint                  `form:"product_id" json:"product_id" binding:"required"`
	Rating        uint                  `form:"rating" json:"rating" binding:"required"`
	Review        string                `form:"review" json:"review" binding:"required"`
	Photo1        *multipart.FileHeader `form:"photo_1"`
	Photo2        *multipart.FileHeader `form:"photo_2"`
	Photo3        *multipart.FileHeader `form:"photo_3"`
	Photo4        *multipart.FileHeader `form:"photo_4"`
	Photo5        *multipart.FileHeader `form:"photo_5"`
	Video         *multipart.FileHeader `form:"video"`
}

type ProductReviewCreateResponse struct {
	ID            uint
	UserID        uuid.UUID
	OrderDetailID uuid.UUID
	ProductID     uint
	Rating        uint
	Review        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ProductReviewGetByProductIDResponse struct {
	ID                 uint
	UserID             uuid.UUID
	User               models.User
	ProductID          uint      `json:"product_id"`
	Rating             uint      `json:"rating"`
	Review             string    `json:"review"`
	ProductReviewMedia []*string `json:"product_review_media"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
