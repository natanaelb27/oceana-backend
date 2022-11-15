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
	result := r.db.Where("id = ?", &newModel.ID).Updates(&newModel)
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
