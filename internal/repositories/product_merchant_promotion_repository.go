package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductMerchantPromotionRepository interface {
	CreateProductMerchantPromotion(newModels []*models.ProductMerchantPromotion) (models []*models.ProductMerchantPromotion, err error)
	GetByMerchantPromotionID(merchantPromotionID uint) (models []*models.ProductMerchantPromotion, err error)
	Update(newModel *models.ProductMerchantPromotion) (model *models.ProductMerchantPromotion, err error)
	DeleteByProductID(productID uint) (model *models.ProductMerchantPromotion, err error)
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
