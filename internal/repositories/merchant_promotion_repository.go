package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type MerchantPromotionRepository interface {
	CreateMyShopPromotion(merchantPromotion *models.MerchantPromotion) (*models.MerchantPromotion, error)
}

type merchantPromotionRepository struct {
	db *gorm.DB
}

type MerchantPromotionRepositoryConfig struct {
	DB *gorm.DB
}

func NewMerchantPromotionRepository(c *MerchantPromotionRepositoryConfig) MerchantPromotionRepository {
	return &merchantPromotionRepository{
		db: c.DB,
	}
}

func (mp *merchantPromotionRepository) CreateMyShopPromotion(newModel *models.MerchantPromotion) (model *models.MerchantPromotion, err error) {
	queryResult := mp.db.Create(&newModel)
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
