package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type WalletTopUpRepository interface {
	Create(newModel *models.WalletTopUp) (model *models.WalletTopUp, err error)
}

type walletTopUpRepository struct {
	db *gorm.DB
}

type WalletTopUpRepositoryConfig struct {
	DB *gorm.DB
}

func NewWalletTopupRepository(c *WalletTopUpRepositoryConfig) WalletTopUpRepository {
	return &walletTopUpRepository{
		db: c.DB,
	}
}

func (w *walletTopUpRepository) Create(newModel *models.WalletTopUp) (model *models.WalletTopUp, err error) {
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
