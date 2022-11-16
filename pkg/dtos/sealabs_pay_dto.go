package dtos

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type SealabsPayCreateRequest struct {
	AccountNumber string `json:"account_number" binding:"required"`
}

func (source *SealabsPayCreateRequest) ToModel(id uuid.UUID) (target *models.SealabsPay) {
	target = &models.SealabsPay{
		UserID:        id,
		AccountNumber: source.AccountNumber,
	}
	return
}

type SealabsPayCreateResponse struct {
	ID            uuid.UUID `json:"id"`
	AccountNumber string    `json:"account_number"`
}

func NewSealabsPayCreateResponse(source *models.SealabsPay) (target *SealabsPayCreateResponse) {
	target = &SealabsPayCreateResponse{
		ID:            source.ID,
		AccountNumber: source.AccountNumber,
	}
	return
}
