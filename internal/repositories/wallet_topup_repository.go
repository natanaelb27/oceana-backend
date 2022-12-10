package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type WalletTopupRepository interface {
	Create(newModel *models.WalletTopup) (model *models.WalletTopup, err error)
}

type walletTopupRepository struct {
	db *gorm.DB
}

type WalletTopupRepositoryConfig struct {
	DB *gorm.DB
}

func NewWalletTopupRepository(c *WalletTopupRepositoryConfig) WalletTopupRepository {
	return &walletTopupRepository{
		db: c.DB,
	}
}

func (w *walletTopupRepository) Create(newModel *models.WalletTopup) (model *models.WalletTopup, err error) {
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
