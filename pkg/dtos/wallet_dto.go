package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
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

type WalletDetailsResponse struct {
	ID                       uuid.UUID                        `json:"id"`
	UserID                   uuid.UUID                        `json:"user_id"`
	Balance                  uint                             `json:"balance"`
	WalletTransactionHistory []*dtos.WalletTransactionHistory `json:"wallet_transaction_history"`
	IsActive                 bool                             `json:"is_active"`
	InvalidPinCount          int                              `json:"invalid_pin_count"`
	BlockedExpirationTime    *time.Time                       `json:"blocked_expiration_time,omitempty"`
	CreatedAt                time.Time                        `json:"created_at"`
	UpdatedAt                time.Time                        `json:"updated_at"`
}

func NewWalletDetailsResponse(sourceModel *models.Wallet, sourceDto []*dtos.WalletTransactionHistory) (target *WalletDetailsResponse) {
	target = &WalletDetailsResponse{
		ID:                       sourceModel.ID,
		UserID:                   sourceModel.UserID,
		Balance:                  sourceModel.Balance,
		WalletTransactionHistory: sourceDto,
		IsActive:                 sourceModel.IsActive,
		InvalidPinCount:          sourceModel.InvalidPinCount,
		BlockedExpirationTime:    sourceModel.BlockedExpirationTime,
		CreatedAt:                sourceModel.CreatedAt,
		UpdatedAt:                sourceModel.UpdatedAt,
	}
	return
}
