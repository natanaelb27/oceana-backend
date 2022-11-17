package dtos

import (
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

type AddressGetResponse struct {
	ID            uuid.UUID           `json:"id"`
	UserID        uuid.UUID           `json:"user_id"`
	User          UserCreateResponse  `json:"user"`
	Name          string              `json:"name"`
	Detail        string              `json:"detail"`
	SubDistrictId uint                `json:"sub_district_id"`
	SubDistrict   SubDistrictResponse `json:"sub_district"`
	IsDefault     bool                `json:"is_default"`
}

func NewAddressGetResponse(source *models.Address) (target *AddressGetResponse) {
	target = &AddressGetResponse{
		ID:     source.ID,
		UserID: source.UserID,
		User: UserCreateResponse{
			ID:          source.User.ID,
			FullName:    source.User.FullName,
			Username:    source.User.Username,
			Email:       source.User.Email.Address,
			PhoneNumber: source.User.PhoneNumber,
			Gender:      source.User.Gender,
			Birthdate:   source.User.Birthdate,
			Role:        source.User.Role,
		},
		Name:          source.Name,
		Detail:        source.Detail,
		SubDistrictId: source.SubDistrictID,
		SubDistrict: SubDistrictResponse{
			ID:         source.SubDistrict.ID,
			Name:       source.SubDistrict.Name,
			DistrictID: source.SubDistrict.DistrictID,
			District: DistrictResponse{
				ID:     source.SubDistrict.District.ID,
				Name:   source.SubDistrict.District.Name,
				CityID: source.SubDistrict.District.CityID,
				City: CityResponse{
					ID:         source.SubDistrict.District.City.ID,
					Name:       source.SubDistrict.District.City.Name,
					ProvinceID: source.SubDistrict.District.City.ProvinceID,
					Province: ProvinceResponse{
						ID:   source.SubDistrict.District.City.Province.ID,
						Name: source.SubDistrict.District.City.Province.Name,
					},
				},
			},
			ZipCode: source.SubDistrict.ZipCode,
		},
		IsDefault: source.IsDefault,
	}

	return
}
