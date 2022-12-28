package dtos

import (
	"time"

	"github.com/google/uuid"
)

type MerchantInformation struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Photo           *string   `json:"photo"`
	TotalRating     uint      `json:"total_rating"`
	RatingStar      float64   `json:"rating_star"`
	ProductCount    uint      `json:"product_count"`
	UserID          uuid.UUID `json:"user_id"`
	ShopAddress     string    `json:"shop_address"`
	ShopSubDistrict string    `json:"shop_sub_district"`
	ShopDistrict    string    `json:"shop_district"`
	ShopCity        string    `json:"shop_city"`
	ShopProvince    string    `json:"shop_province"`
	JoinedDate      time.Time `json:"joined_date"`
}
