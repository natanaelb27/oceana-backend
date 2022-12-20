package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type OrderDetailList struct {
	OrderDetailID         uuid.UUID      `json:"order_detail_id"`
	Message               string         `json:"message"`
	ProductImage          string         `json:"product_image"`
	ProductName           string         `json:"product_name"`
	ProductVariations     string         `json:"product_variations"`
	Quantity              uint           `json:"quantity"`
	Price                 uint           `json:"price"`
	TotalPrice            uint           `json:"total_price"`
	ProductReviewID       *uint          `json:"product_review_id"`
	Rating                *uint          `json:"rating"`
	Review                *string        `json:"review"`
	ReviewImages          pq.StringArray `json:"review_images" gorm:"type:text[]"`
	ReviewVideos          pq.StringArray `json:"review_videos" gorm:"type:text[]"`
	ReviewCreatedDatetime *time.Time     `json:"review_created_datetime"`
}
