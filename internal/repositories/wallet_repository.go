package repositories

import (
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
