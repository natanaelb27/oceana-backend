package repositories

import (
	"database/sql"
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	internalProductDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/product_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductItemRepository interface {
	WithTx(tx *gorm.DB) *productItemRepository
	GetByProductID(productID uint) ([]*dtos.ProductModel, error)
	GetById(id uint) (model *models.ProductItem, err error)
	CreateProductItem(productItem *models.ProductItem) (model *models.ProductItem, err error)
	GetDetailsById(ID uint) (*dtos.ProductItemDetails, error)
	GetUpdateByProductID(productID uint) ([]*internalProductDtos.ProductItem, error)
	SubtractProductStock(ID uint, amount uint) (*models.ProductItem, error)
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

func (u productItemRepository) WithTx(tx *gorm.DB) *productItemRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *productItemRepository) GetByProductID(productID uint) ([]*dtos.ProductModel, error) {
	var productModels []*dtos.ProductModel

	subQueryVariationIndex := r.db.Model(&models.ProductVariant{}).
		Select(`
			product_variants.product_variant_type_id AS product_variant_type_id, 
			product_variants.id AS product_variant_id, 
			ROW_NUMBER () OVER (PARTITION BY product_variants.product_variant_type_id ORDER BY product_variants.id ASC) - 1 AS variation_index,
			name AS name
		`)

	subQueryProductItemVariationIndex := r.db.Model(&models.ProductItem{}).
		Select(`
			product_items.id AS product_item_id, 
			STRING_AGG(tbl_variation_index.name, ', ' ORDER BY tbl_variation_index.product_variant_type_id ASC) AS product_variations,
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
			product_item_variation_index.variation_index AS variation_index,
			product_item_variation_index.product_variations AS product_variations
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

func (r *productItemRepository) GetDetailsById(ID uint) (*dtos.ProductItemDetails, error) {
	var productItemDetails *dtos.ProductItemDetails
	result := r.db.Model(&models.ProductItem{}).
		Select(`
			product_items.id AS product_item_id,
			product_items.sku AS sku,
			product_items.stock_quantity AS stock_quantity,
			product_items.price AS price,
			products.category_id AS category_id,
			products.id AS product_id,
			products.merchant_id AS merchant_id,
			products.name AS name,
			products.description AS description,
			products.is_hazardous AS is_hazardous,
			products.weight AS weight,
			products.package_size AS package_size,
			products.condition AS condition
		`).
		Joins("JOIN products ON products.id = product_items.product_id").
		Where("product_items.id = ?", ID).
		First(&productItemDetails)

	if result.RowsAffected == 0 {
		return nil, e.ErrDataNotFound
	}

	return productItemDetails, nil
}

func (r *productItemRepository) GetUpdateByProductID(productID uint) ([]*internalProductDtos.ProductItem, error) {
	var productModels []*internalProductDtos.ProductItem

	sqVariant := r.db.Table("product_variants p_variants").
		Select(`
			p_variants.id as id,
			-- get row number (start from 0 because of -1 in the end) of product variants partitioned by product_variants_type_id that which ordered ascending by  product variants id
			(ROW_NUMBER() OVER (PARTITION BY p_variants.product_variant_type_id ORDER BY p_variants.id ASC))-1 AS index
			`)
	sqVariantIndex := r.db.Model(&models.ProductVariation{}).
		Select("ARRAY_AGG(p_variants.index)").
		Joins("JOIN (?) p_variants ON p_variants.id = product_variations.product_variant_id", sqVariant).
		Where("product_variations.product_item_id = pi.id")
	queryResult := r.db.Table("product_items pi").Debug().
		Select(`
			pi.id AS id,
			pi.sku AS sku,
			pi.stock_quantity AS stock_quantity,
			pi.price AS price,
			COALESCE((@sqVariantIndex),'{}') AS variant_index 
			`,
			sql.Named("sqVariantIndex", sqVariantIndex),
		).
		Joins("JOIN products ON pi.product_id = products.id").
		Where("products.id = ?", productID).
		Order("pi.id ASC").
		Find(&productModels)
	if queryResult.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	return productModels, nil
}

func (r *productItemRepository) SubtractProductStock(ID uint, amount uint) (*models.ProductItem, error) {
	var model *models.ProductItem

	result := r.db.
		Model(&models.ProductItem{}).
		Where("id = ?", ID).
		Update("stock_quantity", gorm.Expr("stock_quantity - ?", amount))

	if result.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrDataNotFound
	}

	return model, nil
}
