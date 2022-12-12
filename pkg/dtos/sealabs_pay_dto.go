package dtos

import (
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type SealabsPayCreateRequest struct {
	NameOnCard string `json:"name_on_card" binding:"required"`
	CardNumber string `json:"card_number" binding:"required"`
	TxnID      string `json:"txn_id" binding:"required"`
	Signature  string `json:"signature" binding:"required"`
}

func (source *SealabsPayCreateRequest) ToModel(id uuid.UUID) (target *models.SealabsPay) {
	target = &models.SealabsPay{
		UserID:        id,
		AccountNumber: source.CardNumber,
		NameOnCard:    source.NameOnCard,
		ActiveDate:    time.Now(),
	}
	return
}

type SealabsPayCreateResponse struct {
	ID         uuid.UUID `json:"id"`
	CardNumber string    `json:"card_number"`
	ActiveDate time.Time `json:"active_date"`
	NameOnCard string    `json:"name_on_card"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewSealabsPayCreateResponse(source *models.SealabsPay) (target *SealabsPayCreateResponse) {
	target = &SealabsPayCreateResponse{
		ID:         source.ID,
		CardNumber: source.AccountNumber,
		ActiveDate: source.ActiveDate,
		NameOnCard: source.NameOnCard,
		CreatedAt:  source.CreatedAt,
		UpdatedAt:  source.UpdatedAt,
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
	NameOnCard string `json:"name_on_card"`
	CardNumber string `json:"card_number" binding:"required"`
	Amount     string `json:"amount" binding:"required"`
}

type SealabsPayPaymentResponseError struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type SealabsPayRefundRequest struct {
	TxnID  int    `json:"txn_id"`
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
}

func NewSealabsPayRefundRequest(txnID string) (*SealabsPayRefundRequest, error) {
	txnIDInt, err := strconv.Atoi(txnID)
	if err != nil {
		return nil, e.ErrInternalOperationFailed
	}
	amountInt, err := strconv.Atoi(constants.SealabsPayRegistrationAmount)
	if err != nil {
		return nil, e.ErrInternalOperationFailed
	}
	return &SealabsPayRefundRequest{
		TxnID:  txnIDInt,
		Amount: amountInt,
		Reason: constants.SealabsPayRegistrationReason,
	}, nil
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
