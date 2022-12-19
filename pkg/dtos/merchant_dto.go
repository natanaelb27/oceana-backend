package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type MerchantCreateRequest struct {
	UserID    uuid.UUID
	Name      string    `json:"name" binding:"required"`
	AddressID uuid.UUID `json:"address_id" binding:"required"`
}

func (source *MerchantCreateRequest) ToModel() (target *models.Merchant) {
	target = &models.Merchant{
		UserID:    source.UserID,
		Name:      source.Name,
		AddressID: source.AddressID,
	}
	return
}

type MerchantCreateResponse struct {
	ID        uint      `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	AddressID uuid.UUID `json:"address_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewMerchantCreateResponse(source *models.Merchant) (target *MerchantCreateResponse) {
	target = &MerchantCreateResponse{
		ID:        source.ID,
		UserID:    source.UserID,
		Name:      source.Name,
		AddressID: source.AddressID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
	return
}

type MerchantGetInformationByIDResponse struct {
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

func NewMerchantGetInformationByIDResponse(source *dtos.MerchantInformation) (target *MerchantGetInformationByIDResponse) {
	target = &MerchantGetInformationByIDResponse{
		ID:              source.ID,
		Name:            source.Name,
		Photo:           source.Photo,
		TotalRating:     source.TotalRating,
		RatingStar:      source.RatingStar,
		ProductCount:    source.ProductCount,
		UserID:          source.UserID,
		ShopAddress:     source.ShopAddress,
		ShopSubDistrict: source.ShopSubDistrict,
		ShopDistrict:    source.ShopDistrict,
		ShopCity:        source.ShopCity,
		ShopProvince:    source.ShopProvince,
		JoinedDate:      source.JoinedDate,
	}
	return
}

type MerchantGetProfileResponse struct {
	ID        uint      `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	AddressID uuid.UUID `json:"address_id"`
	Photo     *string   `json:"photo"`
}

func NewMerchantGetProfileResponse(source *models.Merchant) (target *MerchantGetProfileResponse) {
	target = &MerchantGetProfileResponse{
		ID:        source.ID,
		UserID:    source.UserID,
		Name:      source.Name,
		AddressID: source.AddressID,
		Photo:     source.Photo,
	}
	return
}
