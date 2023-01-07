package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/dtos"
)

type PendingPaymentGetAllByUserResponse struct {
	TransactionID             uuid.UUID `json:"transaction_id"`
	UserID                    uuid.UUID `json:"user_id"`
	UserFullName              string    `json:"user_full_name"`
	TransactionDatetime       time.Time `json:"transaction_datetime"`
	PaymentMethodID           uint      `json:"payment_method_id"`
	PaymentMethodName         string    `json:"payment_method_name"`
	PaymentMethodImage        string    `json:"payment_method_image"`
	PaymentExpirationDatetime time.Time `json:"payment_expiration_datetime"`
	CardNumber                string    `json:"card_number"`
	TotalPayment              float64   `json:"total_payment"`
}

func NewPendingPaymentGetAllByUserResponse(source []*dtos.PendingPaymentTransaction) (target []*PendingPaymentGetAllByUserResponse) {
	for _, result := range source {
		target = append(target, &PendingPaymentGetAllByUserResponse{
			TransactionID:             result.TransactionID,
			UserID:                    result.UserID,
			UserFullName:              result.UserFullName,
			TransactionDatetime:       result.TransactionDatetime,
			PaymentMethodID:           result.PaymentMethodID,
			PaymentMethodName:         result.PaymentMethodName,
			PaymentMethodImage:        result.PaymentMethodImage,
			PaymentExpirationDatetime: result.PaymentExpirationDatetime,
			CardNumber:                result.CardNumber,
			TotalPayment:              result.TotalPayment,
		})
	}

	return
}
