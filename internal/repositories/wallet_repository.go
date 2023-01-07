package repositories

import (
	"database/sql"
	"errors"
	"log"
	"math"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository interface {
	WithTx(tx *gorm.DB) *walletRepository
	GetByUserId(id uuid.UUID) (model *models.Wallet, err error)
	Activate(id uuid.UUID, newModel *models.Wallet) (model *models.Wallet, err error)
	Create(newModel *models.Wallet) (model *models.Wallet, err error)
	Update(id uuid.UUID, newModel *models.Wallet) (model *models.Wallet, err error)
	GetWalletTransactionHistory(id uuid.UUID, pagination *utils.Pagination, query *externalDtos.WalletDetailsQuery) ([]*dtos.WalletTransactionHistory, *utils.Pagination, error)
}

type walletRepository struct {
	db *gorm.DB
}

type WalletRepositoryConfig struct {
	DB *gorm.DB
}

func NewWalletRepository(c *WalletRepositoryConfig) WalletRepository {
	return &walletRepository{
		db: c.DB,
	}
}

func (u walletRepository) WithTx(tx *gorm.DB) *walletRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *walletRepository) GetByUserId(id uuid.UUID) (model *models.Wallet, err error) {
	result := r.db.Where("user_id = ?", id).First(&model)
	if result.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *walletRepository) Activate(id uuid.UUID, newModel *models.Wallet) (model *models.Wallet, err error) {
	result := r.db.Where("id = ?", id).Updates(&newModel)
	if result.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	return
}

func (r *walletRepository) Update(id uuid.UUID, newModel *models.Wallet) (model *models.Wallet, err error) {
	result := r.db.Where("id = ?", id).Updates(&newModel)
	if result.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	return
}

func (w *walletRepository) Create(newModel *models.Wallet) (model *models.Wallet, err error) {
	queryResult := w.db.Create(&newModel)
	if queryResult.RowsAffected == 0 {
		err = e.ErrInsertFailed
		return
	}
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}
	model = newModel
	return
}

func (w *walletRepository) GetWalletTransactionHistory(id uuid.UUID, pagination *utils.Pagination, query *externalDtos.WalletDetailsQuery) ([]*dtos.WalletTransactionHistory, *utils.Pagination, error) {
	var walletTransactionHistory []*dtos.WalletTransactionHistory
	var wallet *models.Wallet
	var merchant *models.Merchant
	result := w.db.Model(&models.Wallet{}).First(&wallet, id)
	if result.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}
	if result.RowsAffected == 0 {
		return nil, nil, e.ErrDataNotFound
	}
	queryWalletTopup := w.db.Model(&models.WalletTopup{}).
		Select(`
			id, 
			wallet_id, 
			amount, 
			transaction_date AS transaction_datetime, 
			'Top Up Transaction' AS transaction_type,
			'Credit' AS type
		`).
		Where("wallet_id = ?", id)

	queryWalletTransactionsAsBuyer := w.db.Model(&models.Transaction{}).
		Select(`
			transactions.id AS id, 
			wallet_transaction_payments.wallet_id AS wallet_id, 
			transactions.total_price_with_shipping_fee AS amount, 
			transactions.transaction_datetime, 
			'Payment Transaction' AS transaction_type,
			'Debit' AS type
		`).
		Joins("JOIN wallet_transaction_payments ON transactions.id = wallet_transaction_payments.transaction_id").
		Where("wallet_transaction_payments.wallet_id = ? AND transactions.user_id = ?", id, wallet.UserID)

	subQueryRefund := w.db.Model(&models.Refund{}).
		Select(`
			refunds.id AS refund_id,
			SUM(refund_details.amount_user_get) AS amount
		`).
		Joins("JOIN refund_details ON refunds.id = refund_details.refund_id").
		Where("refunds.status = ?", constants.RefundStatusCompleted).
		Group("refunds.id")

	queryRefundTransactions := w.db.Model(&models.Refund{}).
		Select(`
				refunds.id AS id,
				wallet_transaction_payments.wallet_id,
				(
					CASE
						WHEN orders.status = @canceled_by_system THEN subQueryRefund.amount + orders.shipping_fee
						WHEN orders.status = @canceled_by_seller THEN subQueryRefund.amount + orders.shipping_fee
						ELSE subQueryRefund.amount
					END
				) AS amount, 
				refunds.completed_datetime AS transaction_datetime, 
				(
					CASE
						WHEN orders.status = @canceled_by_system THEN 'Canceled Transaction'
						WHEN orders.status = @canceled_by_seller THEN 'Canceled Transaction'
						ELSE 'Refund Transaction'
					END
				) AS transaction_type,
				'Credit' AS type
			`,
			sql.Named("canceled_by_system", constants.SystemCanceledStatus),
			sql.Named("canceled_by_seller", constants.SellerCanceledStatus)).
		Joins("JOIN (?) subQueryRefund ON refunds.id = subQueryRefund.refund_id", subQueryRefund).
		Joins("JOIN orders ON refunds.order_id = orders.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN payment_methods ON transactions.payment_method_id = payment_methods.id").
		Joins("JOIN wallet_transaction_payments ON transactions.id = wallet_transaction_payments.transaction_id").
		Where("payment_methods.id = ? AND wallet_transaction_payments.wallet_id = ?", constants.WalletID, id)

	result = w.db.Model(&models.Merchant{}).Where("user_id = ?", wallet.UserID).First(&merchant)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		subQuery := w.db.Model(&models.RefundDetail{}).Select(`
			SUM(refund_details.amount_merchant_lost)
		`).
			Joins("JOIN refunds ON refund_details.refund_id = refunds.id").
			Where("refunds.status = ?", constants.RefundStatusCompleted).
			Where("refunds.order_id = orders.id").
			Group("refunds.order_id")

		queryWalletTransactionsAsSeller := w.db.Model(&models.Order{}).
			Select(`
				orders.id AS id, 
				wallets.id AS wallet_id, 
				(orders.total_price_after_discount - COALESCE((?), 0)) AS amount,
				orders.claim_datetime AS transaction_datetime, 
				'Received Payment Transaction' AS transaction_type,
				'Credit' AS type
			`, subQuery).
			Joins("JOIN transactions ON transactions.id = orders.transaction_id").
			Joins("JOIN merchants ON merchants.id = orders.merchant_id").
			Joins("JOIN wallets ON wallets.user_id = merchants.user_id").
			Where("orders.merchant_id = ?", merchant.ID).
			Where("orders.status = ?", constants.CompletedStatus).
			Where("orders.claim_status = true")

		queryResult := w.db.Raw("(?) UNION (?) UNION (?) UNION (?)",
			queryWalletTopup,
			queryWalletTransactionsAsBuyer,
			queryRefundTransactions,
			queryWalletTransactionsAsSeller,
		)

		if query.FilterType != "" && query.FilterType != "all" {
			queryResult = w.db.Table("(?) AS tbl", queryResult).Where("type ILIKE ?", query.FilterType)
		}

		if query.FilterStartTransactionDate != nil && query.FilterEndTransactionDate != nil {
			queryResult = w.db.Table("(?) AS tbl", queryResult).Where("transaction_datetime::date BETWEEN ?::date AND ?::date", query.FilterStartTransactionDate, query.FilterEndTransactionDate)
		}

		queryResult = w.db.Table("(?) AS tbl", queryResult).Order("transaction_datetime DESC")

		var totalRows int64
		queryResult = w.db.Table("(?) AS tbl", queryResult).Count(&totalRows)

		if queryResult.Error != nil {
			return nil, nil, e.ErrDatabaseOperationFailed
		}

		pagination.TotalRows = uint(totalRows)
		pagination.TotalPages = uint(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
		queryResult = w.db.Table("(?) AS tbl", queryResult).Offset(int(pagination.GetOffset())).Limit(int(pagination.GetLimit())).Find(&walletTransactionHistory)

		if queryResult.Error != nil {
			return nil, nil, e.ErrDatabaseOperationFailed
		}

	} else {
		queryResult := w.db.Raw("(?) UNION (?) UNION (?)",
			queryWalletTopup,
			queryWalletTransactionsAsBuyer,
			queryRefundTransactions,
		)

		if query.FilterType != "" && query.FilterType != "all" {
			queryResult = w.db.Table("(?) AS tbl", queryResult).Where("type ILIKE ?", query.FilterType)
		}

		if query.FilterStartTransactionDate != nil && query.FilterEndTransactionDate != nil {
			queryResult = w.db.Table("(?) AS tbl", queryResult).Where("transaction_datetime::date BETWEEN ?::date AND ?::date", query.FilterStartTransactionDate, query.FilterEndTransactionDate)
		}

		queryResult = w.db.Table("(?) AS tbl", queryResult).Order("transaction_datetime DESC")

		var totalRows int64
		queryResult = w.db.Table("(?) AS tbl", queryResult).Count(&totalRows)

		if queryResult.Error != nil {
			return nil, nil, e.ErrDatabaseOperationFailed
		}

		pagination.TotalRows = uint(totalRows)
		pagination.TotalPages = uint(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
		queryResult = w.db.Table("(?) AS tbl", queryResult).Offset(int(pagination.GetOffset())).Limit(int(pagination.GetLimit())).Find(&walletTransactionHistory)

		if queryResult.Error != nil {
			return nil, nil, e.ErrDatabaseOperationFailed
		}
	}

	return walletTransactionHistory, pagination, nil
}
