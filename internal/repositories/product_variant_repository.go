package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductVariantRepository interface {
	Create(productVariant *models.ProductVariant) (model *models.ProductVariant, err error)
	GetAll() (model []*models.ProductVariant, err error)
	GetById(id uint) (model *models.ProductVariant, err error)
	Update(id uint, newModel *models.ProductVariant) (model *models.ProductVariant, err error)
	Delete(id uint) (model *models.ProductVariant, err error)
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

func (r *productVariantRepository) Delete(id uint) (model *models.ProductVariant, err error) {
	result := r.db.Delete(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model.ID = id

	return
}
