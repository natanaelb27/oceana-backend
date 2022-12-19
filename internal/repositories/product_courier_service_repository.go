package repositories

import (
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductCourierServiceRepository interface {
	WithTx(tx *gorm.DB) *productCourierServiceRepository
	Create(productCourierService *models.ProductCourierService) (model *models.ProductCourierService, err error)
	GetAll() (model []*models.ProductCourierService, err error)
	GetById(id uint) (model *models.ProductCourierService, err error)
	GetByConditionMap(conditionMap map[string]any) (model *models.ProductCourierService, err error)
	Update(id uint, newModel *models.ProductCourierService) (model *models.ProductCourierService, err error)
	Delete(id uint) (model *models.ProductCourierService, err error)
	DeleteByMerchantCourierServiceIDProductID(merchantCourierServiceID uint, productID uint) (err error)
}

type productCourierServiceRepository struct {
	db *gorm.DB
}

type ProductCourierServiceRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductCourierServiceRepository(c *ProductCourierServiceRepositoryConfig) ProductCourierServiceRepository {
	return &productCourierServiceRepository{
		db: c.DB,
	}
}

func (u productCourierServiceRepository) WithTx(tx *gorm.DB) *productCourierServiceRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *productCourierServiceRepository) Create(newModel *models.ProductCourierService) (model *models.ProductCourierService, err error) {
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

func (r *productCourierServiceRepository) GetAll() (model []*models.ProductCourierService, err error) {
	queryResult := r.db.Find(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *productCourierServiceRepository) GetById(id uint) (model *models.ProductCourierService, err error) {
	result := r.db.First(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *productCourierServiceRepository) GetByConditionMap(conditionMap map[string]any) (model *models.ProductCourierService, err error) {
	result := r.db.Where(conditionMap).First(&model)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *productCourierServiceRepository) Update(id uint, newModel *models.ProductCourierService) (model *models.ProductCourierService, err error) {
	result := r.db.Model(&model).Where("id = ?", id).Updates(&newModel)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	model.ID = id

	return
}

func (r *productCourierServiceRepository) Delete(id uint) (model *models.ProductCourierService, err error) {
	result := r.db.Delete(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model.ID = id

	return
}

func (r *productCourierServiceRepository) DeleteByMerchantCourierServiceIDProductID(merchantCourierServiceID uint, productID uint) (err error) {
	result := r.db.
		Where("merchant_courier_service_id = ?", merchantCourierServiceID).
		Where("product_id = ?", productID).
		Delete(&models.ProductCourierService{})
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}
