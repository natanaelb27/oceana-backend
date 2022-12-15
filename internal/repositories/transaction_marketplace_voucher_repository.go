package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type TransactionMarketplaceVoucherRepository interface {
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

func (r *transactionMarketplaceVoucherRepository) Create(newModel *models.TransactionMarketplaceVoucher) (model *models.TransactionMarketplaceVoucher, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}
