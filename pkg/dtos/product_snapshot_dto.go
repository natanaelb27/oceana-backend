package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type GetByOrderDetailIDResponse struct {
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

func NewGetByOrderDetailIDResponse(source *dtos.ProductSnapshotDetails) (target *GetByOrderDetailIDResponse) {
	target = &GetByOrderDetailIDResponse{
		ProductSnapshotID:   source.ProductSnapshotID,
		UserID:              source.UserID,
		TransactionDatetime: source.TransactionDatetime,
		ProductID:           source.ProductID,
		MerchantPhoto:       source.MerchantPhoto,
		MerchantID:          source.MerchantID,
		ShopName:            source.ShopName,
		Name:                source.Name,
		IsHazardous:         source.IsHazardous,
		Weight:              source.Weight,
		PackageSize:         source.PackageSize,
		Condition:           source.Condition,
		Description:         source.Description,
		CurrentStatus:       source.CurrentStatus,
		Images:              source.Images,
		Videos:              source.Videos,
		ProductVariations:   source.ProductVariations,
		Price:               source.Price,
	}

	return
}
