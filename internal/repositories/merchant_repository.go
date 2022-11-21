package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	Create(merchant *models.Merchant) (*models.Merchant, error)
	GetByName(name string) (*models.Merchant, error)
	GetByUserID(id uuid.UUID) (model *models.Merchant, err error)
	UpdateProfileAddressShopDefault(userID uuid.UUID, addressID uuid.UUID) (model *models.Merchant, err error)
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

func (m *merchantRepository) UpdateProfileAddressShopDefault(userID uuid.UUID, addressID uuid.UUID) (model *models.Merchant, err error) {

	result := m.db.Model(&model).Where("user_id = ?", userID).Update("address_id", addressID)

	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	getUpdated := m.db.Where("address_id = ?", addressID).Find(&model)

	if getUpdated.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}
