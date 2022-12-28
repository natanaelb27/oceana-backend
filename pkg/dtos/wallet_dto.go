package dtos

import (
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
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

type WalletDetailsQuery struct {
	FilterType                 string     `json:"filter-type"`
	FilterStartTransactionDate *time.Time `json:"filter-start-transaction-date"`
	FilterEndTransactionDate   *time.Time `json:"filter-end-transaction-date"`
	Limit                      uint       `json:"limit"`
	Page                       uint       `json:"page"`
}

func NewWalletDetailsQuery(c *gin.Context) *WalletDetailsQuery {
	var query WalletDetailsQuery

	filterType := c.Query("filter-type")
	filterStartTransactionDate := c.Query("filter-start-transaction-date")
	filterEndTransactionDate := c.Query("filter-end-transaction-date")
	limit := c.Query("limit")
	page := c.Query("page")

	query.FilterType = filterType

	if filterStartTransactionDate != "" && filterEndTransactionDate != "" {
		startDate, err := time.Parse("2006-01-02", filterStartTransactionDate)

		if err != nil {
			query.FilterStartTransactionDate = nil
		} else {
			query.FilterStartTransactionDate = &startDate
		}

		endDate, err := time.Parse("2006-01-02", filterEndTransactionDate)

		if err != nil {
			query.FilterEndTransactionDate = nil
		} else {
			query.FilterEndTransactionDate = &endDate
		}

	} else {
		query.FilterStartTransactionDate = nil
		query.FilterEndTransactionDate = nil
	}

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		query.Page = 1
	} else if pageInt <= 0 {
		query.Page = 1
	} else {
		query.Page = uint(pageInt)
	}

	limitInt, err := strconv.Atoi(limit)

	if err != nil {
		query.Limit = 10
	} else if limitInt <= 0 {
		query.Limit = 10
	} else {
		query.Limit = uint(limitInt)
	}

	return &query
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

type TopupWalletRequest struct {
	CardNumber   string `json:"card_number"`
	TxID         string `json:"txn_id"`
	Amount       string `json:"amount"`
	Merchantcode string `json:"merchant_code"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	Signature    string `json:"signature"`
}

func (source *TopupWalletRequest) ToModel() (target *models.WalletTopup, err error) {
	amount, err := strconv.Atoi(source.Amount)
	if err != nil {
		err = e.ErrInvalidAmount
		return
	}

	target = &models.WalletTopup{
		TxID:            source.TxID,
		Amount:          uint(amount),
		TransactionDate: time.Now(),
	}

	return
}
