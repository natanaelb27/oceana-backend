package repositories

import (
	"log"

	internalDtos "github.com/natanaelb27/oceana-backend/internal/dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductVariantRepository interface {
	WithTx(tx *gorm.DB) *productVariantRepository
	Create(productVariant *models.ProductVariant) (model *models.ProductVariant, err error)
	GetAll() (model []*models.ProductVariant, err error)
	GetById(id uint) (model *models.ProductVariant, err error)
	GetAllByProductVariantTypeID(productVariantTypeID uint) (model []*models.ProductVariant, err error)
	Update(id uint, newModel *models.ProductVariant) (model *models.ProductVariant, err error)
	Delete(id uint) (model *models.ProductVariant, err error)
	GetAllVariantIndexByProductID(productID uint) (productModel []*internalDtos.ProductVariationTypeListOfVariants, err error)
}

type productVariantRepository struct {
	db *gorm.DB
}

type ProductVariantRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductVariantRepository(c *ProductVariantRepositoryConfig) ProductVariantRepository {
	return &productVariantRepository{
		db: c.DB,
	}
}

func (u productVariantRepository) WithTx(tx *gorm.DB) *productVariantRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *productVariantRepository) Create(newModel *models.ProductVariant) (model *models.ProductVariant, err error) {
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

func (r *productVariantRepository) GetAll() (model []*models.ProductVariant, err error) {
	queryResult := r.db.Find(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *productVariantRepository) GetById(id uint) (model *models.ProductVariant, err error) {
	result := r.db.First(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *productVariantRepository) GetAllByProductVariantTypeID(productVariantTypeID uint) (model []*models.ProductVariant, err error) {
	result := r.db.
		Where("product_variant_type_id = ?", productVariantTypeID).
		Order("id ASC").
		Find(&model)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *productVariantRepository) Update(id uint, newModel *models.ProductVariant) (model *models.ProductVariant, err error) {
	result := r.db.Model(&model).Where("id = ?", id).Updates(&newModel)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	model.ID = id

	return
}

func (r *productVariantRepository) UpdateBy(id uint, newModel *models.ProductVariant) (model *models.ProductVariant, err error) {
	result := r.db.Model(&model).Where("id = ?", id).Updates(&newModel)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	model.ID = id

	return
}

func (r *productVariantRepository) Delete(id uint) (model *models.ProductVariant, err error) {
	result := r.db.Delete(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model.ID = id

	return
}

func (r *productVariantRepository) GetAllVariantIndexByProductID(productID uint) (productModel []*internalDtos.ProductVariationTypeListOfVariants, err error) {
	sqProductVariants := r.db.Model(&models.ProductVariant{}).
		Select(`
			ARRAY_AGG(product_variants.id ORDER BY product_variants.id) AS agg,
			product_variants.product_variant_type_id AS product_variant_type_id
		`).
		Joins("JOIN product_variant_types ON product_variants.product_variant_type_id = product_variant_types.id").
		Where("product_variants.product_variant_type_id = pvt.id").
		Group("product_variants.product_variant_type_id")
	result := r.db.Debug().
		Table("product_variant_types pvt").
		Select(`
			pv.agg AS variants
		`,
		).
		Joins("JOIN LATERAL (?) pv ON pvt.id = pv.product_variant_type_id", sqProductVariants).
		Where("pvt.product_id = ?", productID).
		Find(&productModel)
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
