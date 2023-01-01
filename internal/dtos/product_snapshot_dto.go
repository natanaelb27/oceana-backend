package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ProductSnapshotDetails struct {
	ProductSnapshotID   uint           `json:"product_snapshot_id"`
	UserID              uuid.UUID      `json:"user_id"`
	TransactionDatetime time.Time      `json:"transaction_datetime"`
	ProductID           uint           `json:"product_id"`
	MerchantPhoto       string         `json:"merchant_photo"`
	MerchantID          uint           `json:"merchant_id"`
	ShopName            string         `json:"shop_name"`
	Name                string         `json:"name"`
	IsHazardous         bool           `json:"is_hazardous"`
	Weight              uint           `json:"weight"`
	PackageSize         string         `json:"package_size"`
	Condition           string         `json:"condition"`
	Description         string         `json:"description"`
	CurrentStatus       string         `json:"current_status"`
	Images              pq.StringArray `json:"images" gorm:"type:text[]"`
	Videos              pq.StringArray `json:"videos" gorm:"type:text[]"`
	ProductVariations   string         `json:"product_variations"`
	Price               float64        `json:"price"`
}
