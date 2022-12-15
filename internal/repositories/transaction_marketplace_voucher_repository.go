package repositories

import (
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type TransactionMarketplaceVoucherRepository interface {
	WithTx(tx *gorm.DB) *transactionMarketplaceVoucherRepository
	Create(newModel *models.TransactionMarketplaceVoucher) (model *models.TransactionMarketplaceVoucher, err error)
}

type transactionMarketplaceVoucherRepository struct {
	db *gorm.DB
}

type TransactionMarketplaceVoucherRepositoryConfig struct {
	DB *gorm.DB
}

func NewTransactionMarketplaceVoucherRepository(c *TransactionMarketplaceVoucherRepositoryConfig) TransactionMarketplaceVoucherRepository {
	return &transactionMarketplaceVoucherRepository{
		db: c.DB,
	}
}

func (u transactionMarketplaceVoucherRepository) WithTx(tx *gorm.DB) *transactionMarketplaceVoucherRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *transactionMarketplaceVoucherRepository) Create(newModel *models.TransactionMarketplaceVoucher) (model *models.TransactionMarketplaceVoucher, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}
