package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/models"
)

type SealabsPayDTO struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	CardNumber string
	ActiveDate time.Time
	NameOnCard string
}

func NewSealabsPayDTO(source *models.SealabsPay) (target *SealabsPayDTO) {
	target = &SealabsPayDTO{
		ID:         source.ID,
		UserID:     source.UserID,
		CardNumber: source.AccountNumber,
		ActiveDate: source.ActiveDate,
		NameOnCard: source.NameOnCard,
	}
	return
}

type SealabsPayInfo struct {
	ApiURL       string
	ApiKey       string
	MerchantCode string
	RedirectURL  string
	CallbackURL  string
	Signature    string
}
