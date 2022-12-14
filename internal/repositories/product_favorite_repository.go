package repositories

import (
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductFavoriteRepository interface {
	WithTx(tx *gorm.DB) *productFavoriteRepository
	Create(productFavorite *models.ProductFavorite) (model *models.ProductFavorite, err error)
	GetByUserIDAndProductID(userID uuid.UUID, productID uint) (model *models.ProductFavorite, err error)
	Update(id uint, newModel *models.ProductFavorite) (model *models.ProductFavorite, err error)
}

type productFavoriteRepository struct {
	db *gorm.DB
}

type ProductFavoriteRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductFavoriteRepository(c *ProductFavoriteRepositoryConfig) ProductFavoriteRepository {
	return &productFavoriteRepository{
		db: c.DB,
	}
}

func (u productFavoriteRepository) WithTx(tx *gorm.DB) *productFavoriteRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *productFavoriteRepository) Create(newModel *models.ProductFavorite) (model *models.ProductFavorite, err error) {
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

func (r *productFavoriteRepository) GetByUserIDAndProductID(userID uuid.UUID, productID uint) (model *models.ProductFavorite, err error) {
	result := r.db.First(&model, "user_id = ? AND product_id = ?", userID, productID)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *productFavoriteRepository) Update(id uint, newModel *models.ProductFavorite) (model *models.ProductFavorite, err error) {
	result := r.db.Model(&model).Where("id = ?", id).Update("is_favorite", newModel.IsFavorite)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	model.ID = id

	return
}
