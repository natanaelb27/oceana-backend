package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductItemRepository interface {
	GetByProductID(productID uint) ([]*dtos.ProductModel, error)
	GetById(id uint) (model *models.ProductItem, err error)
	CreateProductItem(productItem *models.ProductItem) (model *models.ProductItem, err error)
}

type productItemRepository struct {
	db *gorm.DB
}

type ProductItemRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductItemRepository(c *ProductItemRepositoryConfig) ProductItemRepository {
	return &productItemRepository{
		db: c.DB,
	}
}

func (r *productItemRepository) GetByProductID(productID uint) ([]*dtos.ProductModel, error) {
	var productModels []*dtos.ProductModel

	subQueryVariationIndex := r.db.Model(&models.ProductVariant{}).
		Select(`
		product_variants.product_variant_type_id AS product_variant_type_id, 
		product_variants.id AS product_variant_id, 
		ROW_NUMBER () OVER (PARTITION BY product_variants.product_variant_type_id ORDER BY product_variants.id ASC) - 1 AS variation_index 
		`)

	subQueryProductItemVariationIndex := r.db.Model(&models.ProductItem{}).
		Select(`
			product_items.id AS product_item_id, 
			ARRAY_AGG(variation_index ORDER BY tbl_variation_index.product_variant_type_id ASC) AS variation_index
		`).
		Joins("JOIN product_variations ON product_variations.product_item_id = product_items.id").
		Joins("JOIN (?) tbl_variation_index ON tbl_variation_index.product_variant_id = product_variations.product_variant_id", subQueryVariationIndex).
		Group("product_items.id")

	subQueryProductDiscount := r.db.Model(&models.ProductMerchantPromotion{}).
		Select(`
			product_merchant_promotions.product_id,
			product_merchant_promotions.discount_type,
			product_merchant_promotions.discount,
			product_merchant_promotions.maximum_discount_amount
		`).
		Joins("JOIN merchant_promotions ON merchant_promotions.id = product_merchant_promotions.merchant_promotion_id").
		Where(`
			product_merchant_promotions.maximum_quota > 
			CASE
				WHEN product_merchant_promotions.maximum_quota > -1 THEN product_merchant_promotions.current_quota
				ELSE -2
			END
		`).
		Where("NOW() BETWEEN merchant_promotions.start_datetime AND merchant_promotions.end_datetime").
		Where("product_merchant_promotions.product_id = ?", productID)

	queryResult := r.db.Model(&models.ProductItem{}).
		Select(`
			product_items.id AS product_item_id,
			product_items.sku AS sku,
			product_items.stock_quantity AS stock_quantity,
			product_items.price AS price,
			(
				CASE 
					WHEN LOWER(product_discount_agg.discount_type) = 'nominal' 
					THEN GREATEST(product_items.price - product_discount_agg.discount, 0)::float
					ELSE product_items.price - LEAST((product_items.price * product_discount_agg.discount) / 100.0, product_discount_agg.maximum_discount_amount)
				END
			) AS price_after_discount,
			product_item_variation_index.variation_index AS variation_index
		`).
		Joins("JOIN products ON product_items.product_id = products.id").
		Joins("LEFT JOIN (?) product_discount_agg ON products.id = product_discount_agg.product_id", subQueryProductDiscount).
		Joins("LEFT JOIN (?) product_item_variation_index ON product_items.id = product_item_variation_index.product_item_id", subQueryProductItemVariationIndex).
		Where("products.id = ?", productID).
		Find(&productModels)

	if queryResult.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	return productModels, nil
}

func (r *productItemRepository) GetById(id uint) (model *models.ProductItem, err error) {
	result := r.db.First(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (pi *productItemRepository) CreateProductItem(newModel *models.ProductItem) (model *models.ProductItem, err error) {
	queryResult := pi.db.Create(&newModel)
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
