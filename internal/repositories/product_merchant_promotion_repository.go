package repositories

import (
	"fmt"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductMerchantPromotionRepository interface {
	CreateProductMerchantPromotion(productMerchantPromotion []*models.ProductMerchantPromotion) ([]*models.ProductMerchantPromotion, error)
}

type productMerchantPromotionRepository struct {
	db *gorm.DB
}

type ProductMerchantPromotionRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductMerchantPromotionRepositoryConfig(c *ProductMerchantPromotionRepositoryConfig) ProductMerchantPromotionRepository {
	return &productMerchantPromotionRepository{
		db: c.DB,
	}
}

func (pmp *productMerchantPromotionRepository) CreateProductMerchantPromotion(newModels []*models.ProductMerchantPromotion) (models []*models.ProductMerchantPromotion, err error) {
	for _, newModel := range newModels {
		fmt.Println(newModel.ID)
	}
	queryResult := pmp.db.Create(&newModels)
	if queryResult.RowsAffected == 0 {
		err = e.ErrInsertFailed
		return
	}
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	models = newModels

	return
}
