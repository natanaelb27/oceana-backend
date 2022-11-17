package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type WalletActivateRequest struct {
	Pin             string `json:"pin" binding:"required,len=6"`
	PinConfirmation string `json:"pin_confirmation" binding:"required,eqfield=Pin"`
}

func (source *WalletActivateRequest) ToModel() (target *models.Wallet) {
	target = &models.Wallet{
		Pin:      source.Pin,
		IsActive: true,
	}
	return
}

type WalletActivateResponse struct {
	ID        uuid.UUID `json:"id"`
	Balance   uint      `json:"balance"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewWalletActivateResponse(source *models.Wallet) (target *WalletActivateResponse) {
	target = &WalletActivateResponse{
		ID:        source.ID,
		Balance:   source.Balance,
		IsActive:  source.IsActive,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
	return
}

type WalletChangePinRequest struct {
	UpgradeToken    string `json:"upgrade_token" binding:"required"`
	Pin             string `json:"pin" binding:"required,len=6"`
	PinConfirmation string `json:"pin_confirmation" binding:"required,eqfield=Pin"`
}

func (source *WalletChangePinRequest) ToModel() (target *models.Wallet) {
	target = &models.Wallet{
		Pin: source.Pin,
	}
	return
}

type WalletChangePinResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewWalletChangePinResponse(source *models.Wallet) (target *WalletChangePinResponse) {
	target = &WalletChangePinResponse{
		ID:        source.ID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
	return
}
