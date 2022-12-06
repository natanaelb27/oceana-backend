package repositories

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductMerchantPromotionRepository interface {
	CreateProductMerchantPromotion(productMerchantPromotion []*models.ProductMerchantPromotion) ([]*models.ProductMerchantPromotion, error)
	GetOnGoingByProductID(productID uint) (*models.ProductMerchantPromotion, error)
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

func (pmp *productMerchantPromotionRepository) GetOnGoingByProductID(productID uint) (*models.ProductMerchantPromotion, error) {
	var model *models.ProductMerchantPromotion

	result := pmp.db.
		Preload("MerchantPromotion").
		Joins("LEFT JOIN merchant_promotions ON merchant_promotion_id = merchant_promotions.id").
		First(&model, "product_id = ? AND merchant_promotions.start_datetime >= ? AND merchant_promotions.end_datetime <= ?", productID, time.Now(), time.Now())

	if result.Error == gorm.ErrRecordNotFound {
		return nil, e.ErrProductNotFound
	}

	return model, result.Error
}
