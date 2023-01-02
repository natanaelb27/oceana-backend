package dtos

import (
	"time"
)

type MarketplaceVoucherDetails struct {
	ID                    uint      `json:"id"`
	BannerUrl             string    `json:"banner_url"`
	Name                  string    `json:"name"`
	StartDatetime         time.Time `json:"start_datetime"`
	EndDatetime           time.Time `json:"end_datetime"`
	DiscountType          string    `json:"discount_type"`
	Discount              uint      `json:"discount"`
	MinimumSpendAmount    uint      `json:"minimum_spend_amount"`
	MaximumDiscountAmount uint      `json:"maximum_discount_amount"`
	MaximumQuota          uint      `json:"maximum_quota"`
	CurrentQuota          uint      `json:"current_quota"`
	VoucherStatus         string    `json:"voucher_status"`
}
