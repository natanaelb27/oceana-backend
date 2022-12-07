package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductVariantMediaRepository interface {
	Create(productVariantMedia *models.ProductVariantMedia) (model *models.ProductVariantMedia, err error)
	GetAll() (model []*models.ProductVariantMedia, err error)
	GetById(id uint) (model *models.ProductVariantMedia, err error)
	Update(id uint, newModel *models.ProductVariantMedia) (model *models.ProductVariantMedia, err error)
	Delete(id uint) (model *models.ProductVariantMedia, err error)
}

type productVariantMediaRepository struct {
	db *gorm.DB
}

type ProductVariantMediaRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductVariantMediaRepository(c *ProductVariantMediaRepositoryConfig) ProductVariantMediaRepository {
	return &productVariantMediaRepository{
		db: c.DB,
	}
}

func (r *productVariantMediaRepository) Create(newModel *models.ProductVariantMedia) (model *models.ProductVariantMedia, err error) {
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

func (r *productVariantMediaRepository) GetAll() (model []*models.ProductVariantMedia, err error) {
	queryResult := r.db.Find(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *productVariantMediaRepository) GetById(id uint) (model *models.ProductVariantMedia, err error) {
	result := r.db.First(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *productVariantMediaRepository) Update(id uint, newModel *models.ProductVariantMedia) (model *models.ProductVariantMedia, err error) {
	result := r.db.Model(&model).Where("id = ?", id).Updates(&newModel)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	model.ID = id

	return
}

func (r *productVariantMediaRepository) Delete(id uint) (model *models.ProductVariantMedia, err error) {
	result := r.db.Delete(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model.ID = id

	return
}
