package dtos

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type AddressCreateRequest struct {
	UserID        uuid.UUID `json:"user_id" binding:"required"`
	Name          string    `form:"name" json:"name" binding:"required"`
	Detail        string    `form:"detail" json:"detail" binding:"required"`
	SubDistrictID uint      `form:"sub_district_id" json:"sub_district_id" binding:"required"`
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
