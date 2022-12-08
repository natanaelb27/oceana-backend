package dtos

import "github.com/google/uuid"

type AddressDetails struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	PhoneNumber   string    `json:"phone_number"`
	Name          string    `json:"name"`
	Detail        string    `json:"detail"`
	SubDistrictID uint      `json:"sub_district_id"`
	SubDistrict   string    `json:"sub_district"`
	ZipCode       string    `json:"zip_code"`
	DistrictID    uint      `json:"district_id"`
	District      string    `json:"district"`
	CityID        uint      `json:"city_id"`
	City          string    `json:"city"`
	ProvinceID    uint      `json:"province_id"`
	Province      string    `json:"province"`
	IsDefault     bool      `json:"is_default"`
	IsShopAddress bool      `json:"is_shop_address"`
}
