package repositories

import (
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MerchantConfigurationRepository interface {
	WithTx(tx *gorm.DB) *merchantConfigurationRepository
	GetByMerchantID(id uint) (model *models.MerchantConfiguration, err error)
	Create(newModel *models.MerchantConfiguration) (model *models.MerchantConfiguration, err error)
	UpsertUsingMap(modelMap map[string]any) (err error)
}

type merchantConfigurationRepository struct {
	db *gorm.DB
}

type MerchantConfigurationRepositoryConfig struct {
	DB *gorm.DB
}

func NewMerchantConfigurationRepository(c *MerchantConfigurationRepositoryConfig) MerchantConfigurationRepository {
	return &merchantConfigurationRepository{
		db: c.DB,
	}
}

func (u merchantConfigurationRepository) WithTx(tx *gorm.DB) *merchantConfigurationRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *merchantConfigurationRepository) GetByMerchantID(id uint) (model *models.MerchantConfiguration, err error) {
	result := r.db.Where("merchant_id = ?", id).First(&model)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if result.Error == gorm.ErrRecordNotFound {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *merchantConfigurationRepository) UpsertUsingMap(modelMap map[string]any) (err error) {
	result := r.db.
		Model(&models.MerchantConfiguration{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "merchant_id"}},
			DoUpdates: clause.Assignments(modelMap),
		}).Create(modelMap)
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

func (w *merchantConfigurationRepository) Create(newModel *models.MerchantConfiguration) (model *models.MerchantConfiguration, err error) {
	queryResult := w.db.Create(&newModel)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if queryResult.RowsAffected == 0 {
		err = e.ErrInsertFailed
		return
	}

	model = newModel

	return
}
