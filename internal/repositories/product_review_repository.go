package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductReviewRepository interface {
	Create(newModel *models.ProductReview) (model *models.ProductReview, err error)
	GetByProductID(id uint) ([]*models.ProductReview, error)
}

type productReviewRepository struct {
	db *gorm.DB
}

type ProductReviewRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductReviewRepository(c *ProductReviewRepositoryConfig) ProductReviewRepository {
	return &productReviewRepository{
		db: c.DB,
	}
}

func (pr *productReviewRepository) Create(newModel *models.ProductReview) (model *models.ProductReview, err error) {
	queryResult := pr.db.Create(&newModel)
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

func (pr *productReviewRepository) GetByProductID(id uint) ([]*models.ProductReview, error) {
	var productReviews []*models.ProductReview 
	result := pr.db.Model(&models.ProductReview{}).Where("product_id = ?", id).Preload("User").Find(&productReviews)
	if result.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}
	return productReviews, nil
}
