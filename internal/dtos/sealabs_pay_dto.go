package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type SealabsPayDTO struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	AccountNumber string
	ActiveDate    time.Time
	NameOnCard    string
}

func NewSealabsPayDTO(source *models.SealabsPay) (target *SealabsPayDTO) {
	target = &SealabsPayDTO{
		ID:            source.ID,
		UserID:        source.UserID,
		AccountNumber: source.AccountNumber,
		ActiveDate:    source.ActiveDate,
		NameOnCard:    source.NameOnCard,
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
