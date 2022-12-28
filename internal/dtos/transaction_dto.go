package dtos

import (
	"time"

	"github.com/google/uuid"
)

type PendingPaymentTransaction struct {
	TransactionID             uuid.UUID `json:"transaction_id"`
	UserID                    uuid.UUID `json:"user_id"`
	UserFullName              string    `json:"user_full_name"`
	TransactionDatetime       time.Time `json:"transaction_datetime"`
	PaymentMethodID           uint      `json:"payment_method_id"`
	PaymentMethodName         string    `json:"payment_method_name"`
	PaymentMethodImage        string    `json:"payment_method_image"`
	PaymentExpirationDatetime time.Time `json:"payment_expiration_datetime"`
	CardNumber                string    `json:"card_number" gorm:"-"`
	TotalPayment              float64   `json:"total_payment"`
}
