package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductItemRepository interface {
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
