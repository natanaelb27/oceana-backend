package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	Create(merchant *models.Merchant) (*models.Merchant, error)
	GetByName(name string) (*models.Merchant, error) 
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

func (m *merchantRepository) Create(newModel *models.Merchant) (model *models.Merchant, err error) {
	queryResult := m.db.Create(&newModel)
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

func (m *merchantRepository) GetByName(name string) (model *models.Merchant, err error) {
	result := m.db.Model(&models.Merchant{}).Where("name = ?", name).Find(&model)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}
	return
}
