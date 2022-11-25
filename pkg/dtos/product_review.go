package dtos

import (
	"mime/multipart"
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/gin-gonic/gin"
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

type ProductReviewQuery struct {
	Limit uint `json:"limit"`
	Page  uint `json:"page"`
}

func NewProductReviewQuery(c *gin.Context) *ProductReviewQuery {
	var query *ProductReviewQuery

	limit := c.Query("limit")
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)

	if err != nil {
		query.Page = 1
	} else if pageInt <= 0 {
		query.Page = 1
	} else {
		query.Page = uint(pageInt)
	}

	limitInt, err := strconv.Atoi(limit)

	if err != nil {
		query.Limit = 10
	} else if limitInt <= 0 {
		query.Limit = 10
	} else {
		query.Limit = uint(limitInt)
	}

	return query
}

type ProductReviewGetByProductIDResponse struct {
	Pagination    *utils.Pagination              `json:"pagination"`
	ProductReview []*ProductReviewGetByProductID `json:"product_review"`
}

func NewProductReviewGetByIDResponse(source []*ProductReviewGetByProductID, pagination *utils.Pagination) (target *ProductReviewGetByProductIDResponse) {
	if pagination == nil {
		return nil
	}

	target = &ProductReviewGetByProductIDResponse{
		Pagination:    pagination,
		ProductReview: source,
	}

	return
}
