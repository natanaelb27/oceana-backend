package repositories

import (
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	internalProductDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/product_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductVariationRepository interface {
	WithTx(tx *gorm.DB) *productVariationRepository
	GetByProductID(productID uint) ([]*dtos.ProductVariation, error)
	GetUpdateByProductID(productID uint) ([]*internalProductDtos.Variants, error)
	Create(newModel *models.ProductVariation) (model *models.ProductVariation, err error)
	DeleteAllByProductItemID(id uint) (err error)
}

type productVariationRepository struct {
	db *gorm.DB
}

type ProductVariationRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductVariationRepository(c *ProductVariationRepositoryConfig) ProductVariationRepository {
	return &productVariationRepository{
		db: c.DB,
	}
}

func (u productVariationRepository) WithTx(tx *gorm.DB) *productVariationRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *productVariationRepository) GetByProductID(productID uint) ([]*dtos.ProductVariation, error) {
	var productVariations []*dtos.ProductVariation

	subQueryOptions := r.db.Model(&models.ProductVariant{}).
		Select(`
			product_variant_types.id AS product_variant_type_id, 
			ARRAY_AGG(product_variants.name ORDER BY product_variants.id ASC) AS options
		`).
		Joins("JOIN product_variant_types ON product_variant_types.id = product_variants.product_variant_type_id").
		Group("product_variant_types.id").
		Order("product_variant_types.id ASC")

	subQuerySumStockQuantity := r.db.Model(&models.ProductItem{}).
		Select(`
			product_variants.product_variant_type_id AS product_variant_type_id, 
			product_variants.id AS product_variant_id, 
			SUM(product_items.stock_quantity) AS stock_quantity
		`).
		Joins("JOIN product_variations ON product_variations.product_item_id = product_items.id").
		Joins("JOIN product_variants ON product_variants.id = product_variations.product_variant_id").
		Group("product_items.product_id, product_variants.product_variant_type_id, product_variants.id").
		Order("product_variants.product_variant_type_id ASC, product_variants.id ASC")

	subQuerySummedStock := r.db.Table("(?) AS sum_stock_quantity", subQuerySumStockQuantity).
		Select(`
			sum_stock_quantity.product_variant_type_id,
			ARRAY_AGG(sum_stock_quantity.stock_quantity ORDER BY sum_stock_quantity.product_variant_id ASC) AS summed_stocks
		`).
		Group("product_variant_type_id")

	subQueryVariantImage := r.db.Model(&models.ProductVariantMedia{}).
		Select(`
			product_variants.product_variant_type_id AS product_variant_type_id, 
			ARRAY_AGG(product_variant_media.media_url ORDER BY product_variant_media.product_variant_id ASC) AS images
		`).
		Joins("JOIN product_variants ON product_variant_media.product_variant_id = product_variants.id").
		Group("product_variants.product_variant_type_id")

	queryResult := r.db.Model(&models.ProductVariantType{}).
		Select(`
			product_variant_types.name AS name,
			option.options AS options,
			COALESCE(variant_image.images, '{}') AS image,
			summed_stock.summed_stocks AS summed_stocks
		`).
		Joins("JOIN (?) option ON product_variant_types.id = option.product_variant_type_id", subQueryOptions).
		Joins("JOIN (?) summed_stock ON product_variant_types.id = summed_stock.product_variant_type_id", subQuerySummedStock).
		Joins("LEFT JOIN (?) variant_image ON product_variant_types.id = variant_image.product_variant_type_id", subQueryVariantImage).
		Where("product_variant_types.product_id = ?", productID).
		Find(&productVariations)

	if queryResult.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	return productVariations, nil
}

func (r *productVariationRepository) GetUpdateByProductID(productID uint) ([]*internalProductDtos.Variants, error) {
	var result []*internalProductDtos.Variants

	subQueryOptions := r.db.Model(&models.ProductVariant{}).
		Select(`
		product_variant_types.id AS product_variant_type_id, 
		ARRAY_AGG(product_variants.name ORDER BY product_variants.id ASC) AS options
	`).
		Joins("JOIN product_variant_types ON product_variant_types.id = product_variants.product_variant_type_id").
		Group("product_variant_types.id").
		Order("product_variant_types.id ASC")
	subQueryVariantImage := r.db.Model(&models.ProductVariantMedia{}).
		Select(`
		product_variants.product_variant_type_id AS product_variant_type_id, 
		ARRAY_AGG(product_variant_media.media_url ORDER BY product_variant_media.product_variant_id ASC) AS images
	`).
		Joins("JOIN product_variants ON product_variant_media.product_variant_id = product_variants.id").
		Group("product_variants.product_variant_type_id")
	queryResult := r.db.Model(&models.ProductVariantType{}).
		Select(`
			product_variant_types.name AS name,
			option.options AS options,
			COALESCE(variant_image.images, '{}') AS images
		`).
		Joins("LEFT JOIN (?) option ON product_variant_types.id = option.product_variant_type_id", subQueryOptions).
		Joins("LEFT JOIN (?) variant_image ON product_variant_types.id = variant_image.product_variant_type_id", subQueryVariantImage).
		Where("product_variant_types.product_id = ?", productID).
		Find(&result)
	if queryResult.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	return result, nil
}

func (r *productVariationRepository) Create(newModel *models.ProductVariation) (model *models.ProductVariation, err error) {
	queryResult := r.db.Create(&newModel)
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

func (r *productVariationRepository) DeleteAllByProductItemID(id uint) (err error) {
	queryResult := r.db.
		Where("product_item_id = ?", id).
		Delete(&models.ProductVariation{})
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}
