package repositories

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductMerchantPromotionRepository interface {
	CreateProductMerchantPromotion(productMerchantPromotion []*models.ProductMerchantPromotion) ([]*models.ProductMerchantPromotion, error)
	GetDetailsByID(id uint) (*dtos.ProductMerchantPromotionDetails, error)
	GetOnGoingByProductID(productID uint) (*models.ProductMerchantPromotion, error)
	AddCurrentQuota(ID uint, amount uint) (*models.ProductMerchantPromotion, error)
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

func (pmp *productMerchantPromotionRepository) GetDetailsByID(id uint) (*dtos.ProductMerchantPromotionDetails, error) {
	var productMerchantPromotionDetails *dtos.ProductMerchantPromotionDetails

	result := pmp.db.Model(&models.ProductMerchantPromotion{}).
		Select(`
			merchant_promotions.id AS merchant_promotion_id,
			merchant_promotions.merchant_id AS merchant_id,
			merchant_promotions.name AS name,
			merchant_promotions.start_datetime AS start_datetime,
			merchant_promotions.end_datetime AS end_datetime,
			product_merchant_promotions.id AS product_merchant_promotion_id,
			product_merchant_promotions.product_id AS product_id,
			product_merchant_promotions.discount_type AS discount_type,
			product_merchant_promotions.discount AS discount,
			product_merchant_promotions.maximum_discount_amount AS maximum_discount_amount,
			product_merchant_promotions.maximum_quota AS maximum_quota,
			product_merchant_promotions.maximum_user_limit AS maximum_user_limit,
			product_merchant_promotions.current_quota AS current_quota,
			(
				CASE WHEN NOW()::timestamp BETWEEN start_datetime::timestamp 
				AND end_datetime::timestamp THEN 'ongoing' 
				WHEN start_datetime::timestamp > NOW()::timestamp  THEN 'upcoming' 
				WHEN end_datetime::timestamp < NOW()::timestamp  THEN 'ended' END
			) AS promotion_status
		`).
		Joins("JOIN merchant_promotions ON product_merchant_promotions.merchant_promotion_id = merchant_promotions.id").
		Where("product_merchant_promotions.id = ?", id).
		First(&productMerchantPromotionDetails)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrDataNotFound
		}

		return nil, e.ErrDatabaseOperationFailed
	}

	return productMerchantPromotionDetails, nil
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

func (pmp *productMerchantPromotionRepository) AddCurrentQuota(ID uint, amount uint) (*models.ProductMerchantPromotion, error) {
	var model *models.ProductMerchantPromotion

	result := pmp.db.
		Model(&models.ProductMerchantPromotion{}).
		Where("id = ?", ID).
		Update("current_quota", gorm.Expr("current_quota + ?", amount))

	if result.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrDataNotFound
	}

	return model, nil

}
