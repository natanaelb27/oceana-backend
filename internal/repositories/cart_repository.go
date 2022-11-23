package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	Create(cart *models.Cart) (model *models.Cart, err error)
	GetAll() (model []*models.Cart, err error)
	GetByUserIdProductItemId(userId uuid.UUID, productItemId uint) (model *models.Cart, err error)
	Update(id uint, newModel *models.Cart) (model *models.Cart, err error)
	Delete(id uint) (model *models.Cart, err error)
}

type cartRepository struct {
	db *gorm.DB
}

type CartRepositoryConfig struct {
	DB *gorm.DB
}

func NewCartRepository(c *CartRepositoryConfig) CartRepository {
	return &cartRepository{
		db: c.DB,
	}
}

func (r *cartRepository) Create(newModel *models.Cart) (model *models.Cart, err error) {
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

func (r *cartRepository) GetAll() (model []*models.Cart, err error) {
	queryResult := r.db.Find(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *cartRepository) GetByUserIdProductItemId(userId uuid.UUID, productItemId uint) (model *models.Cart, err error) {
	result := r.db.Where("user_id = ?", userId).Where("product_item_id = ?", productItemId).First(&model)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if result.Error == gorm.ErrRecordNotFound {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *cartRepository) Update(id uint, newModel *models.Cart) (model *models.Cart, err error) {
	result := r.db.Model(&model).Where("id = ?", id).Updates(&newModel)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	model.ID = id

	return
}

func (r *cartRepository) Delete(id uint) (model *models.Cart, err error) {
	result := r.db.Delete(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model.ID = id

	return
}
