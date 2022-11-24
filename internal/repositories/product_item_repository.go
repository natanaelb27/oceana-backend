package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductItemRepository interface {
	GetById(id uint) (model *models.ProductItem, err error)
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

func (r *productItemRepository) GetById(id uint) (model *models.ProductItem, err error) {
	result := r.db.First(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}
