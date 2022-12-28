package dtos

import (
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductReviewCreateRequest struct {
	UserID        uuid.UUID
	OrderDetailID uuid.UUID
	Rating        uint                  `form:"rating" json:"rating" binding:"required"`
	Review        string                `form:"review" json:"review"`
	Photo1        *multipart.FileHeader `form:"photo_1"`
	Photo2        *multipart.FileHeader `form:"photo_2"`
	Photo3        *multipart.FileHeader `form:"photo_3"`
	Photo4        *multipart.FileHeader `form:"photo_4"`
	Photo5        *multipart.FileHeader `form:"photo_5"`
	Video         *multipart.FileHeader `form:"video"`
}

func (p *ProductReviewCreateRequest) ToModel() *models.ProductReview {
	return &models.ProductReview{
		UserID:        p.UserID,
		OrderDetailID: p.OrderDetailID,
		Rating:        p.Rating,
		Review:        p.Review,
	}
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

func NewProductReviewCreateResponse(source *models.ProductReview) (target *ProductReviewCreateResponse) {
	target = &ProductReviewCreateResponse{
		ID:            source.ID,
		UserID:        source.UserID,
		OrderDetailID: source.OrderDetailID,
		ProductID:     source.ProductID,
		Rating:        source.Rating,
		Review:        source.Review,
		CreatedAt:     source.CreatedAt,
		UpdatedAt:     source.UpdatedAt,
	}
	return
}

type ProductReviewQuery struct {
	Sort    string `json:"sort"`
	SortBy  string `json:"sort-by"`
	Limit   uint   `json:"limit"`
	Page    uint   `json:"page"`
	Media   bool   `json:"media"`
	Comment bool   `json:"comment"`
	Rating  []uint `json:"rating"`
}

func NewProductReviewQuery(c *gin.Context) *ProductReviewQuery {
	var query ProductReviewQuery

	limit := c.Query("limit")
	page := c.Query("page")
	media := c.Query("media")
	comment := c.Query("comment")

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

	query.Sort = utils.QueryParamSort(c.Query("sort"))

	if strings.ToLower(media) == "true" {
		query.Media = true
	} else {
		query.Media = false
	}

	sortBy := strings.ToLower(c.Query("sort-by"))

	switch sortBy {
	case "latest":
		query.SortBy = `"product_reviews"."created_at"`
	case "highest-rating":
		query.SortBy = "rating"
		query.Sort = "DESC"
	case "lowest-rating":
		query.SortBy = "rating"
		query.Sort = "ASC"
	default:
		query.SortBy = `"product_reviews"."created_at"`
	}

	if strings.ToLower(comment) == "true" {
		query.Comment = true
	} else {
		query.Comment = false
	}

	ratingQuery := c.Query("rating")
	if ratingQuery != "" {
		ratingSlices := strings.Split(ratingQuery, ",")
		ratings := []uint{}

		for _, ratingString := range ratingSlices {
			rating, err := strconv.Atoi(ratingString)

			if err == nil {
				ratings = append(ratings, uint(rating))
			}
		}

		query.Rating = ratings
	}

	return &query
}
