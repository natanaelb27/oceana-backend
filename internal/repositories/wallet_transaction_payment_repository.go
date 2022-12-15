package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type WalletTransactionPaymentRepository interface {
	Create(newModel *models.WalletTransactionPayment) (model *models.WalletTransactionPayment, err error)
}

type walletTransactionPaymentRepository struct {
	db *gorm.DB
}

type WalletTransactionPaymentRepositoryConfig struct {
	DB *gorm.DB
}

func NewWalletTransactionPaymentRepository(c *WalletTransactionPaymentRepositoryConfig) WalletTransactionPaymentRepository {
	return &walletTransactionPaymentRepository{
		db: c.DB,
	}
}

func (r *walletTransactionPaymentRepository) Create(newModel *models.WalletTransactionPayment) (model *models.WalletTransactionPayment, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}
