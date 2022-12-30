package repositories

import (
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type WalletTopupRepository interface {
	WithTx(tx *gorm.DB) *walletTopupRepository
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

func (u walletTopupRepository) WithTx(tx *gorm.DB) *walletTopupRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
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
