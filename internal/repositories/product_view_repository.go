package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductViewRepository interface {
	Create(productView *models.ProductView) (model *models.ProductView, err error)
	GetByUserIDAndProductID(userID uuid.UUID, productID uint) (model *models.ProductView, err error)
	Update(id uint, newModel *models.ProductView) (model *models.ProductView, err error)
}

type productViewRepository struct {
	db *gorm.DB
}

type ProductViewRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductViewRepository(c *ProductViewRepositoryConfig) ProductViewRepository {
	return &productViewRepository{
		db: c.DB,
	}
}

func (r *productViewRepository) Create(newModel *models.ProductView) (model *models.ProductView, err error) {
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

func (r *productViewRepository) GetByUserIDAndProductID(userID uuid.UUID, productID uint) (model *models.ProductView, err error) {
	result := r.db.First(&model, "user_id = ? AND product_id = ?", userID, productID)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *productViewRepository) Update(id uint, newModel *models.ProductView) (model *models.ProductView, err error) {
	result := r.db.Model(&model).Where("id = ?", id).Updates(&newModel)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	model.ID = id

	return
}
