package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
)

type MarketplaceVoucherCreateRequest struct {
	Name                   string    `json:"name" binding:"required"`
	StartDatetime          time.Time `json:"start_datetime" binding:"required"`
	EndDatetime            time.Time `json:"end_datetime" binding:"required"`
	DiscountType           string    `json:"discount_type" binding:"required"`
	Discount               uint      `json:"discount" binding:"required"`
	MinimumSpendAmount     uint      `json:"minimum_spend_amount"`
	MaxmimumDiscountAmount uint      `json:"maximum_spend_amount"`
	MaximumQuota           uint      `json:"maximum_quota" binding:"required"`
}

func (source *MarketplaceVoucherCreateRequest) ToModel() (target *models.MarketplaceVoucher) {
	target = &models.MarketplaceVoucher{
		Name:                   source.Name,
		StartDatetime:          source.StartDatetime,
		EndDatetime:            source.EndDatetime,
		DiscountType:           source.DiscountType,
		Discount:               source.Discount,
		MinimumSpendAmount:     source.MinimumSpendAmount,
		MaxmimumDiscountAmount: source.MaxmimumDiscountAmount,
		MaximumQuota:           source.MaximumQuota,
	}

	return
}

type MarketplaceVoucherCreateResponse struct {
	ID                     uint      `json:"id"`
	Name                   string    `json:"name"`
	StartDatetime          time.Time `json:"start_datetime"`
	EndDatetime            time.Time `json:"end_datetime"`
	DiscountType           string    `json:"discount_type"`
	Discount               uint      `json:"discount"`
	MinimumSpendAmount     uint      `json:"minimum_spend_amount"`
	MaxmimumDiscountAmount uint      `json:"maximum_spend_amount"`
	MaximumQuota           uint      `json:"maximum_quota"`
	CurrentQuota           uint      `json:"current_quota"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

func NewMarketplaceVoucherCreateResponse(source *models.MarketplaceVoucher) (target *MarketplaceVoucherCreateResponse) {
	target = &MarketplaceVoucherCreateResponse{
		ID:                     source.ID,
		Name:                   source.Name,
		StartDatetime:          source.StartDatetime,
		EndDatetime:            source.EndDatetime,
		DiscountType:           source.DiscountType,
		Discount:               source.Discount,
		MinimumSpendAmount:     source.MinimumSpendAmount,
		MaxmimumDiscountAmount: source.MaxmimumDiscountAmount,
		MaximumQuota:           source.MaximumQuota,
		CurrentQuota:           source.CurrentQuota,
		CreatedAt:              source.CreatedAt,
		UpdatedAt:              source.UpdatedAt,
	}

	return
}

type MarketplaceVoucherGetResponse struct {
	ID                     uint      `json:"id"`
	Name                   string    `json:"name"`
	StartDatetime          time.Time `json:"start_datetime"`
	EndDatetime            time.Time `json:"end_datetime"`
	DiscountType           string    `json:"discount_type"`
	Discount               uint      `json:"discount"`
	MinimumSpendAmount     uint      `json:"minimum_spend_amount"`
	MaxmimumDiscountAmount uint      `json:"maximum_spend_amount"`
	MaximumQuota           uint      `json:"maximum_quota"`
	CurrentQuota           uint      `json:"current_quota"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

func NewMarketplaceVoucherGetResponse(source *models.MarketplaceVoucher) (target *MarketplaceVoucherGetResponse) {
	target = &MarketplaceVoucherGetResponse{
		ID:                     source.ID,
		Name:                   source.Name,
		StartDatetime:          source.StartDatetime,
		EndDatetime:            source.EndDatetime,
		DiscountType:           source.DiscountType,
		Discount:               source.Discount,
		MinimumSpendAmount:     source.MinimumSpendAmount,
		MaxmimumDiscountAmount: source.MaxmimumDiscountAmount,
		MaximumQuota:           source.MaximumQuota,
		CurrentQuota:           source.CurrentQuota,
		CreatedAt:              source.CreatedAt,
		UpdatedAt:              source.UpdatedAt,
	}

	return
}

type MarketplaceVoucherDeleteResponse struct {
	ID uint `json:"id"`
}

func NewMarketplaceVoucherDeleteResponse(source *models.MarketplaceVoucher) (target *MarketplaceVoucherDeleteResponse) {
	target = &MarketplaceVoucherDeleteResponse{
		ID: source.ID,
	}

	return
}
