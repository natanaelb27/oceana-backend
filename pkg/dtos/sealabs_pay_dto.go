package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type SealabsPayCreateRequest struct {
	NameOnCard    string `json:"name_on_card" binding:"required"`
	AccountNumber string `json:"account_number" binding:"required"`
}

func (source *SealabsPayCreateRequest) ToModel(id uuid.UUID) (target *models.SealabsPay) {
	target = &models.SealabsPay{
		UserID:        id,
		AccountNumber: source.AccountNumber,
		NameOnCard:    source.NameOnCard,
		ActiveDate:    time.Now(),
	}
	return
}

type SealabsPayCreateResponse struct {
	ID            uuid.UUID `json:"id"`
	AccountNumber string    `json:"account_number"`
	ActiveDate    time.Time `json:"active_date"`
	NameOnCard    string    `json:"name_on_card"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewSealabsPayCreateResponse(source *models.SealabsPay) (target *SealabsPayCreateResponse) {
	target = &SealabsPayCreateResponse{
		ID:            source.ID,
		AccountNumber: source.AccountNumber,
		ActiveDate:    source.ActiveDate,
		NameOnCard:    source.NameOnCard,
		CreatedAt:     source.CreatedAt,
		UpdatedAt:     source.UpdatedAt,
	}
	return
}

type SealabsPayDeleteResponse struct {
	ID uuid.UUID `json:"id"`
}

func NewSealabsPayDeleteResponse(source *models.SealabsPay) (target *SealabsPayDeleteResponse) {
	target = &SealabsPayDeleteResponse{
		ID: source.ID,
	}

	return
}

type SealabsPayPaymentInboundRequest struct {
	CardNumber string `json:"card_number"`
	Amount     string `json:"amount"`
}

type SealabsPayRedirectResponseError struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type SealabsPayRedirectResponseSuccess struct {
	RedirectURL string `json:"redirect_url"`
}

type SealabsPayCallbackRequest struct {
	TxnID        string `json:"txn_id"`
	Amount       string `json:"amount"`
	MerchantCode string `json:"merchant_code"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	Signature    string `json:"signature"`
}
