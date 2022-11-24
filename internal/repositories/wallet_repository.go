package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository interface {
	GetByUserId(id uuid.UUID) (model *models.Wallet, err error)
	Activate(id uuid.UUID, newModel *models.Wallet) (model *models.Wallet, err error)
	Create(newModel *models.Wallet) (model *models.Wallet, err error)
	Update(id uuid.UUID, newModel *models.Wallet) (model *models.Wallet, err error)
	GetWalletTransactionHistory(id uuid.UUID) ([]*dtos.WalletTransactionHistory, error)
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

func (w *walletRepository) GetWalletTransactionHistory(id uuid.UUID) ([]*dtos.WalletTransactionHistory, error) {
	var walletTransactionHistory []*dtos.WalletTransactionHistory
	var wallet *models.Wallet
	var merchant *models.Merchant
	result := w.db.Model(&models.Wallet{}).First(&wallet, id)
	if result.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}
	if result.RowsAffected == 0 {
		return nil, e.ErrDataNotFound
	}
	queryWalletTopUp := w.db.Model(&models.WalletTopup{}).
		Select(`
			id, wallet_id, amount, transaction_date, 'Top Up Transaction' AS transaction_type
		`).
		Where("wallet_id = ?", id)

	queryWalletTransactionsAsBuyer := w.db.Model(&models.Transaction{}).
		Select(`
			transactions.id, wallet_transaction_payments.wallet_id, transactions.total_price_with_shipping_fee, transactions.transaction_datetime, 'Payment Transaction' AS transaction_type
		`).
		Joins("JOIN wallet_transaction_payments ON transactions.id = wallet_transaction_payments.transaction_id").
		Where("wallet_transaction_payments.wallet_id = ? AND transactions.user_id = ?", id, wallet.UserID)

	queryRefundTransactions := w.db.Model(&models.Refund{}).
		Select(`
				transactions.id, wallet_transaction_payments.wallet_id, transactions.total_price_with_shipping_fee, transactions.transaction_datetime, 'Refund Transaction' AS transaction_type
			`).
		Joins("JOIN order_details ON refunds.order_detail_id = order_details.id").
		Joins("JOIN orders ON order_details.order_id = orders.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN payment_methods ON transactions.payment_method_id = payment_methods.id").
		Joins("JOIN wallet_transaction_payments ON transactions.id = wallet_transaction_payments.transaction_id").
		Where("payment_methods.name = ?", "wallet")

	result = w.db.Model(&models.Merchant{}).Where("user_id ?", wallet.UserID).Find(&merchant)
	if result.Error != gorm.ErrRecordNotFound {
		queryWalletTransactionsAsSeller := w.db.Model(&models.Order{}).
			Select(`
				transactions.id, wallet_transaction_payments.wallet_id, orders.total_price_after_discount, transactions.transaction_datetime, 'Received Payment Transaction' AS transaction_type
			`).
			Joins("JOIN transactions ON transactions.id = orders.transaction_id").
			Joins("JOIN wallet_transaction_payments ON transactions.id = wallet_transaction_payments.transaction_id").
			Where("orders.merchant_id = ?", merchant.ID)

		queryResult := w.db.Raw("(?) UNION (?) UNION (?) UNION (?)",
			queryWalletTopUp,
			queryWalletTransactionsAsBuyer,
			queryRefundTransactions,
			queryWalletTransactionsAsSeller,
		).
			Order("transaction_date DESC").
			Find(&walletTransactionHistory)
		if queryResult.Error != nil {
			return nil, e.ErrDatabaseOperationFailed
		}
	} else if result.Error == gorm.ErrRecordNotFound {
		queryResult := w.db.Debug().Raw("(?) UNION (?) UNION (?)",
			queryWalletTopUp,
			queryWalletTransactionsAsBuyer,
			queryRefundTransactions,
		).
			Order("transaction_date DESC").
			Find(&walletTransactionHistory)
		if queryResult.Error != nil {
			return nil, e.ErrDatabaseOperationFailed
		}
	}

	return walletTransactionHistory, nil
}
