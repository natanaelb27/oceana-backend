package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductMediaRepository interface {
	Create(productMedia *models.ProductMedia) (model *models.ProductMedia, err error)
	GetAll() (model []*models.ProductMedia, err error)
	GetById(id uint) (model *models.ProductMedia, err error)
	Update(id uint, newModel *models.ProductMedia) (model *models.ProductMedia, err error)
	Delete(id uint) (model *models.ProductMedia, err error)
}

type productMediaRepository struct {
	db *gorm.DB
}

type ProductMediaRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductMediaRepository(c *ProductMediaRepositoryConfig) ProductMediaRepository {
	return &productMediaRepository{
		db: c.DB,
	}
}

func (r *productMediaRepository) Create(newModel *models.ProductMedia) (model *models.ProductMedia, err error) {
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

func (r *productMediaRepository) GetAll() (model []*models.ProductMedia, err error) {
	queryResult := r.db.Find(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *productMediaRepository) GetById(id uint) (model *models.ProductMedia, err error) {
	result := r.db.First(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *productMediaRepository) Update(id uint, newModel *models.ProductMedia) (model *models.ProductMedia, err error) {
	result := r.db.Model(&model).Where("id = ?", id).Updates(&newModel)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	model.ID = id

	return
}

func (r *productMediaRepository) Delete(id uint) (model *models.ProductMedia, err error) {
	result := r.db.Delete(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model.ID = id

	return
}
