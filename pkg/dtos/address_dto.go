package dtos

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type AddressCreateRequest struct {
	Name          string `form:"name" json:"name" binding:"required"`
	Detail        string `form:"detail" json:"detail" binding:"required"`
	SubDistrictID uint   `form:"sub_district_id" json:"sub_district_id" binding:"required"`
}

type AddressCreateResponse struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Name          string    `json:"name"`
	Detail        string    `json:"detail"`
	SubDistrictID uint      `json:"sub_district_id"`
	IsDefault     bool      `json:"is_default"`
}

func (source *AddressCreateRequest) ToModel(userID uuid.UUID) (target *models.Address) {
	target = &models.Address{
		UserID:        userID,
		Name:          source.Name,
		Detail:        source.Detail,
		SubDistrictID: source.SubDistrictID,
	}

	return
}

func NewAddressCreateResponse(source *models.Address) (target *AddressCreateResponse) {
	target = &AddressCreateResponse{
		ID:            source.ID,
		UserID:        source.UserID,
		Name:          source.Name,
		Detail:        source.Detail,
		SubDistrictID: source.SubDistrictID,
		IsDefault:     source.IsDefault,
	}

	return
}

// type AddressGetResponse struct {
// 	ID          uuid.UUID          `json:"id"`
// 	UserID      uuid.UUID          `json:"user_id"`
// 	User        UserCreateResponse `json:"user"`
// 	Name        string             `json:"name"`
// 	Detail      string             `json:"detail"`
// 	SubDistrict string             `json:"sub_district"`
// 	District    string             `json:"district"`
// 	City        string             `json:"city"`
// 	Province    string             `json:"province"`
// 	ZipCode     string             `json:"zip_code"`
// 	IsDefault   bool               `json:"is_default"`
// }

// func NewAddressGetResponse(source *models.Address) (target *AddressGetResponse) {
// 	target = &AddressGetResponse{
// 		ID:     source.ID,
// 		UserID: source.UserID,
// 		User: UserCreateResponse{
// 			ID:          source.User.ID,
// 			FullName:    source.User.FullName,
// 			Username:    source.User.Username,
// 			Email:       source.User.Email.Address,
// 			PhoneNumber: source.User.PhoneNumber,
// 			Gender:      source.User.Gender,
// 			Birthdate:   source.User.Birthdate,
// 			Role:        source.User.Role,
// 		},
// 		Name:        source.Name,
// 		Detail:      source.Detail,
// 		SubDistrict: source.SubDistrict.Name,
// 		District:    source.SubDistrict.District.Name,
// 		City:        source.SubDistrict.District.City.Name,
// 		Province:    source.SubDistrict.District.City.Province.Name,
// 		ZipCode:     source.SubDistrict.ZipCode,
// 		IsDefault:   source.IsDefault,
// 	}

// 	return
// }

type AddressGetResponse struct {
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

func NewAddressGetResponse(source *dtos.AddressDetails) (target *AddressGetResponse) {
	target = &AddressGetResponse{
		ID:            source.ID,
		UserID:        source.UserID,
		PhoneNumber:   source.PhoneNumber,
		Name:          source.Name,
		Detail:        source.Detail,
		SubDistrictID: source.SubDistrictID,
		SubDistrict:   source.SubDistrict,
		ZipCode:       source.ZipCode,
		DistrictID:    source.DistrictID,
		District:      source.District,
		CityID:        source.CityID,
		City:          source.City,
		ProvinceID:    source.ProvinceID,
		Province:      source.Province,
		IsDefault:     source.IsDefault,
		IsShopAddress: source.IsShopAddress,
	}

	return
}

type AddressDeleteResponse struct {
	ID uuid.UUID `json:"id"`
}

func NewAddressDeleteResponse(source *models.Address) (target *AddressDeleteResponse) {
	target = &AddressDeleteResponse{
		ID: source.ID,
	}

	return
}
