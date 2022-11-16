package dtos

import (
	"time"

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
