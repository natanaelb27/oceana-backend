package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
)

type MerchantVoucherCreateRequest struct {
	Code                  string    `json:"code" binding:"required"`
	StartDatetime         time.Time `json:"start_datetime" binding:"required"`
	EndDatetime           time.Time `json:"end_datetime" binding:"required"`
	DiscountType          string    `json:"discount_type" binding:"required"`
	Discount              uint      `json:"discount" binding:"required"`
	MinimumSpendAmount    uint      `json:"minimum_spend_amount" binding:"required"`
	MaximumDiscountAmount uint      `json:"maximum_discount_amount"`
	MaximumQuota          uint      `json:"maximum_quota" binding:"required"`
	CurrentQuota          uint      `json:"current_quota"`
}

func (source *MerchantVoucherCreateRequest) ToModel(id uint) (target *models.MerchantVoucher) {
	target = &models.MerchantVoucher{
		MerchantID:            id,
		Code:                  source.Code,
		StartDatetime:         source.StartDatetime,
		EndDatetime:           source.EndDatetime,
		DiscountType:          source.DiscountType,
		Discount:              source.Discount,
		MinimumSpendAmount:    source.MinimumSpendAmount,
		MaximumDiscountAmount: source.MaximumDiscountAmount,
		MaximumQuota:          source.MaximumQuota,
		CurrentQuota:          source.CurrentQuota,
	}
	return
}

type MerchantVoucherCreateResponse struct {
	ID                    uint      `json:"id"`
	Code                  string    `json:"code"`
	StartDatetime         time.Time `json:"start_datetime"`
	EndDatetime           time.Time `json:"end_datetime"`
	DiscountType          string    `json:"discount_type"`
	Discount              uint      `json:"discount"`
	MinimumSpendAmount    uint      `json:"minimum_spend_amount"`
	MaximumDiscountAmount uint      `json:"maximum_discount_amount"`
	MaximumQuota          uint      `json:"maximum_quota"`
	CurrentQuota          uint      `json:"current_quota"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func NewMerchantVoucherCreateResponse(source *models.MerchantVoucher) (target *MerchantVoucherCreateResponse) {
	target = &MerchantVoucherCreateResponse{
		ID:                    source.ID,
		Code:                  source.Code,
		StartDatetime:         source.StartDatetime,
		EndDatetime:           source.EndDatetime,
		DiscountType:          source.DiscountType,
		Discount:              source.Discount,
		MinimumSpendAmount:    source.MinimumSpendAmount,
		MaximumDiscountAmount: source.MaximumDiscountAmount,
		MaximumQuota:          source.MaximumQuota,
		CurrentQuota:          source.CurrentQuota,
	}
	return
}
