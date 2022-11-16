package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	GetByUserID(id uuid.UUID) (model *models.Merchant, err error)
}

type merchantRepository struct {
	db *gorm.DB
}

type MerchantRepositoryConfig struct {
	DB *gorm.DB
}

func NewMerchantRepository(c *MerchantRepositoryConfig) MerchantRepository {
	return &merchantRepository{
		db: c.DB,
	}
}

func (r *merchantRepository) GetByUserID(id uuid.UUID) (model *models.Merchant, err error) {
	result := r.db.Where("user_id = ?", id).First(&model)
	if result.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if result.RowsAffected == 0 {
		err = e.ErrMerchantNotFound
		return
	}

	return
}
