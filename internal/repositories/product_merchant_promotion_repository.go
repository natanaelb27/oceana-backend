package repositories

import (
	"log"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants/product_media_constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductMerchantPromotionRepository interface {
	WithTx(tx *gorm.DB) *productMerchantPromotionRepository
	GetDetailsByID(id uint) (*dtos.ProductMerchantPromotionDetails, error)
	GetOnGoingByProductID(productID uint) (*models.ProductMerchantPromotion, error)
	AddCurrentQuota(ID uint, amount uint) (*models.ProductMerchantPromotion, error)
	CreateProductMerchantPromotion(newModels []*models.ProductMerchantPromotion) (models []*models.ProductMerchantPromotion, err error)
	GetByMpAndProductID(merchantPromotionID, productID uint) (model *models.ProductMerchantPromotion, err error)
	GetByMerchantPromotionID(merchantPromotionID uint) (models []*models.ProductMerchantPromotion, err error)
	Update(newModel *models.ProductMerchantPromotion) (model *models.ProductMerchantPromotion, err error)
	DeleteByProductID(productID uint) (model *models.ProductMerchantPromotion, err error)
	GetProductDetailsByMerchantPromotionID(id uint) ([]*dtos.ProductMerchantPromotion, error)
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

func (u productMerchantPromotionRepository) WithTx(tx *gorm.DB) *productMerchantPromotionRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
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

func (pmp *productMerchantPromotionRepository) GetProductDetailsByMerchantPromotionID(id uint) ([]*dtos.ProductMerchantPromotion, error) {
	var productMerchantPromotion []*dtos.ProductMerchantPromotion

	subQueryProductImage := pmp.db.Model(&models.ProductMedia{}).
		Select(`
			product_media.product_id AS product_id,
			product_media.media_url AS media_url
		`).
		Where("product_media.type = ?", product_media_constants.TypeImage).
		Order("product_media.id ASC").
		Limit(1)

	subQueryProductItem := pmp.db.Model(&models.ProductItem{}).
		Select(`
			product_id AS product_id,
			SUM(stock_quantity) AS total_stock,
			MIN(price) AS price_min,
			MAX(price) AS price_max
		`).
		Group("product_items.product_id")

	result := pmp.db.Debug().Model(&models.ProductMerchantPromotion{}).
		Select(`
			product_merchant_promotions.product_id AS product_id,
			products.name AS product_name,
			tbl_product_image.media_url AS product_image,
			product_item_agg.price_max AS price_max,
			(
				CASE 
					WHEN LOWER(product_merchant_promotions.discount_type) = 'nominal' 
					THEN GREATEST(product_item_agg.price_max - product_merchant_promotions.discount, 0)
					ELSE product_item_agg.price_max - LEAST((product_item_agg.price_max * product_merchant_promotions.discount) / 100.0, product_merchant_promotions.maximum_discount_amount)
				END
			) AS price_max_after_discount,
			product_item_agg.price_min AS price_min,
			(
				CASE 
					WHEN LOWER(product_merchant_promotions.discount_type) = 'nominal' 
					THEN GREATEST(product_item_agg.price_min - product_merchant_promotions.discount, 0)
					ELSE product_item_agg.price_min - LEAST((product_item_agg.price_min * product_merchant_promotions.discount) / 100.0, product_merchant_promotions.maximum_discount_amount)
				END
			) AS price_min_after_discount,
			product_item_agg.total_stock AS total_stock,
			product_merchant_promotions.discount_type AS discount_type,
			product_merchant_promotions.discount AS discount,
			product_merchant_promotions.maximum_discount_amount AS maximum_discount_amount,
			product_merchant_promotions.maximum_quota AS maximum_quota,
			product_merchant_promotions.maximum_user_limit AS maximum_user_limit,
			product_merchant_promotions.current_quota AS current_quota
		`).
		Joins("JOIN products ON product_merchant_promotions.product_id = products.id").
		Joins("LEFT JOIN (?) product_item_agg ON products.id = product_item_agg.product_id", subQueryProductItem).
		Joins("LEFT JOIN (?) tbl_product_image ON products.id = tbl_product_image.product_id", subQueryProductImage).
		Where("product_merchant_promotions.merchant_promotion_id = ?", id).
		Find(&productMerchantPromotion)

	if result.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	return productMerchantPromotion, nil
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

func (pmp *productMerchantPromotionRepository) GetByMpAndProductID(merchantPromotionID, productID uint) (model *models.ProductMerchantPromotion, err error) {
	queryResult := pmp.db.Where("merchant_promotion_id = ? AND product_id = ?", merchantPromotionID, productID).First(&model)

	if queryResult.RowsAffected == 0 {
		err = e.ErrProductMerchantPromotionNotFound
		return
	}

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (pmp *productMerchantPromotionRepository) GetByMerchantPromotionID(merchantPromotionID uint) (models []*models.ProductMerchantPromotion, err error) {
	queryResult := pmp.db.Where("merchant_promotion_id = ?", merchantPromotionID).Order("product_id ASC").Find(&models)

	if queryResult.RowsAffected == 0 {
		err = e.ErrProductMerchantPromotionNotFound
		return
	}

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (pmp *productMerchantPromotionRepository) Update(newModel *models.ProductMerchantPromotion) (model *models.ProductMerchantPromotion, err error) {
	queryResult := pmp.db.Where("merchant_promotion_id = ? AND product_id = ?", newModel.MerchantPromotionId, newModel.ProductID).Clauses(clause.Returning{}).Updates(newModel)

	if queryResult.RowsAffected == 0 {
		err = e.ErrProductMerchantPromotionNotFound
		return
	}

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}

func (pmp *productMerchantPromotionRepository) DeleteByProductID(productID uint) (model *models.ProductMerchantPromotion, err error) {
	queryResult := pmp.db.Clauses(clause.Returning{}).Where("product_id = ?", productID).Delete(&model)

	if queryResult.RowsAffected == 0 {
		err = e.ErrProductMerchantPromotionNotFound
		return
	}

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}
