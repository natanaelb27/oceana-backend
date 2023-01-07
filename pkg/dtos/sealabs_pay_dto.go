package dtos

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/constants"
	"github.com/natanaelb27/oceana-backend/internal/models"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
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
	IsDefault  bool      `json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewSealabsPayCreateResponse(source *models.SealabsPay) (target *SealabsPayCreateResponse) {
	target = &SealabsPayCreateResponse{
		ID:         source.ID,
		CardNumber: source.AccountNumber,
		ActiveDate: source.ActiveDate,
		NameOnCard: source.NameOnCard,
		IsDefault:  source.IsDefault,
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
