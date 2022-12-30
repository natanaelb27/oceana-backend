package repositories

import (
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductReviewMediaRepository interface {
	WithTx(tx *gorm.DB) *productReviewMediaRepository
	Create(newModel *models.ProductReviewMedia) (model *models.ProductReviewMedia, err error)
	GetByProductReviewID(id uint) (model []*models.ProductReviewMedia, err error)
}

type productReviewMediaRepository struct {
	db *gorm.DB
}

type ProductReviewMediaRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductReviewMediaRepository(c *ProductReviewMediaRepositoryConfig) ProductReviewMediaRepository {
	return &productReviewMediaRepository{
		db: c.DB,
	}
}

func (u productReviewMediaRepository) WithTx(tx *gorm.DB) *productReviewMediaRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (pr *productReviewMediaRepository) Create(newModel *models.ProductReviewMedia) (model *models.ProductReviewMedia, err error) {
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

func (pr *productReviewMediaRepository) GetByProductReviewID(id uint) (model []*models.ProductReviewMedia, err error) {
	result := pr.db.Where("product_review_id = ?", id).Find(&model)
	if result.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}
	return
}
