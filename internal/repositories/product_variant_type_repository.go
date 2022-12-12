package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductVariantTypeRepository interface {
	Create(productVariantType *models.ProductVariantType) (model *models.ProductVariantType, err error)
	GetAll() (model []*models.ProductVariantType, err error)
	GetById(id uint) (model *models.ProductVariantType, err error)
	Update(id uint, newModel *models.ProductVariantType) (model *models.ProductVariantType, err error)
	Delete(id uint) (model *models.ProductVariantType, err error)
	GetAllByProductID(productID uint) (model []*models.ProductVariantType, err error)
}

type productVariantTypeRepository struct {
	db *gorm.DB
}

type ProductVariantTypeRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductVariantTypeRepository(c *ProductVariantTypeRepositoryConfig) ProductVariantTypeRepository {
	return &productVariantTypeRepository{
		db: c.DB,
	}
}

func (r *productVariantTypeRepository) Create(newModel *models.ProductVariantType) (model *models.ProductVariantType, err error) {
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

func (r *productVariantTypeRepository) GetAll() (model []*models.ProductVariantType, err error) {
	queryResult := r.db.Find(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *productVariantTypeRepository) GetById(id uint) (model *models.ProductVariantType, err error) {
	result := r.db.First(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *productVariantTypeRepository) GetAllByProductID(productID uint) (model []*models.ProductVariantType, err error) {
	result := r.db.Where("product_id = ?", productID).Find(&model)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *productVariantTypeRepository) Update(id uint, newModel *models.ProductVariantType) (model *models.ProductVariantType, err error) {
	queryResult := r.db.Model(&model).Where("id = ?", id).Updates(&newModel)
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel
	model.ID = id

	return
}

func (r *productVariantTypeRepository) Delete(id uint) (model *models.ProductVariantType, err error) {
	result := r.db.Delete(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model.ID = id

	return
}
