package repositories

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	WithTx(tx *gorm.DB) *transactionRepository
	Create(newModel *models.Transaction) (model *models.Transaction, err error)
	PendingPaymentGetAllByUserID(id uuid.UUID) (pendingPaymentTransaction []*dtos.PendingPaymentTransaction, err error)
	GetByID(id uuid.UUID) (model *models.Transaction, err error)
	GetByOrderID(orderID uuid.UUID) (model *models.Transaction, err error)
	Update(model *models.Transaction) (err error)
}

type transactionRepository struct {
	db *gorm.DB
}

type TransactionRepositoryConfig struct {
	DB *gorm.DB
}

func NewTransactionRepository(c *TransactionRepositoryConfig) TransactionRepository {
	return &transactionRepository{
		db: c.DB,
	}
}

func (u transactionRepository) WithTx(tx *gorm.DB) *transactionRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *transactionRepository) Create(newModel *models.Transaction) (model *models.Transaction, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}

func (r *transactionRepository) PendingPaymentGetAllByUserID(id uuid.UUID) (pendingPaymentTransaction []*dtos.PendingPaymentTransaction, err error) {
	paymentExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_PAYMENT_EXPIRATION_INTERVAL"))

	if err != nil {
		paymentExpirationInterval = constants.CronPaymentExpirationInterval
	}

	subQueryPendingPaymentTransactions := r.db.Model(&models.Order{}).
		Select(`
			DISTINCT
			orders.transaction_id AS transaction_id
		`).
		Where("orders.status = ?", constants.PendingPaymentStatus)

	queryResult := r.db.Model(&models.Transaction{}).
		Select(`
			DISTINCT
			transactions.id AS transaction_id,
			transactions.user_id AS user_id,
			users.full_name AS user_full_name,
			transactions.transaction_datetime AS transaction_datetime,
			transactions.payment_method_id AS payment_method_id,
			payment_methods.name AS payment_method_name,
			payment_methods.image AS payment_method_image,
			transactions.transaction_datetime + INTERVAL '1 minute' * @interval_payment_expiration AS payment_expiration_datetime,
			transactions.total_price_with_shipping_fee AS total_payment
		`, sql.Named("interval_payment_expiration", paymentExpirationInterval)).
		Joins("JOIN payment_methods ON transactions.payment_method_id = payment_methods.id").
		Joins("JOIN users ON transactions.user_id = users.id").
		Joins("JOIN (?) tbl_pending_payment ON transactions.id = tbl_pending_payment.transaction_id", subQueryPendingPaymentTransactions).
		// Where("orders.status = ?", constants.PendingPaymentStatus).
		Where("transactions.is_paid = false").
		Where("transactions.user_id = ?", id).
		Find(&pendingPaymentTransaction)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *transactionRepository) GetByID(id uuid.UUID) (model *models.Transaction, err error) {
	queryResult := r.db.First(&model, id)

	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrTransactionNotFound
		}

		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *transactionRepository) GetByOrderID(orderID uuid.UUID) (model *models.Transaction, err error) {
	queryResult := r.db.Model(&models.Transaction{}).
		Joins("JOIN orders ON transactions.id = orders.transaction_id").
		Where("orders.id = ?", orderID).
		First(&model)

	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrOrderNotFound
		}

		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *transactionRepository) Update(model *models.Transaction) (err error) {
	queryResult := r.db.Updates(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}
