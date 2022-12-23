package repositories

import (
	"database/sql"
	"log"

	internalLogisticDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/logistic_dtos"
	internalProductDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/product_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MerchantCourierServiceRepository interface {
	WithTx(tx *gorm.DB) *merchantCourierServiceRepository
	GetAllByCourierID(courierID uint, merchantID uint) (model []*internalLogisticDtos.CourierService, err error)
	GetByMerchantIDCourierServiceID(merchantID uint, courierServiceID uint) (model *models.MerchantCourierService, err error)
	UpsertUsingMap(modelMap map[string]any) (err error)
	Create(newModel *models.MerchantCourierService) (model *models.MerchantCourierService, err error)
	GetAllCouriersByMerchantIDForProductUpdate(merchantID uint) (model []*internalProductDtos.Courier, err error)
	GetAllCourierServicesForProductUpdate(courierID uint, merchantID uint, productID uint) (model []*internalProductDtos.CourierService, err error)
	GetByConditionMap(conditionMap map[string]any) (model *models.MerchantCourierService, err error)
}

type merchantCourierServiceRepository struct {
	db *gorm.DB
}

type MerchantCourierServiceRepositoryConfig struct {
	DB *gorm.DB
}

func NewMerchantCourierServiceRepository(c *MerchantCourierServiceRepositoryConfig) MerchantCourierServiceRepository {
	return &merchantCourierServiceRepository{
		db: c.DB,
	}
}

func (u merchantCourierServiceRepository) WithTx(tx *gorm.DB) *merchantCourierServiceRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *merchantCourierServiceRepository) subQueryMerchantCourierService(merchantID uint) *gorm.DB {
	query := r.db.Table("merchant_courier_services mcs").
		Select(` 
			CASE
				WHEN mcs.is_active = TRUE
					THEN TRUE
				ELSE FALSE
			END
			AS is_selected`).
		Where("mcs.courier_service_id = courier_services.id").
		Where("mcs.is_active = TRUE").
		Where("mcs.merchant_id = ?", merchantID).
		Limit(1)
	return query
}

func (r *merchantCourierServiceRepository) GetAllByCourierID(courierID uint, merchantID uint) (model []*internalLogisticDtos.CourierService, err error) {
	subQueryMerchantCourierService := sql.Named("subQueryMerchantCourierService", r.subQueryMerchantCourierService(merchantID))
	queryResult := r.db.
		Table("courier_services").
		Select(`
			courier_services.id AS id,
			courier_services.code AS code,
			courier_services.description AS description,
			(@subQueryMerchantCourierService) AS is_selected
			`,
			subQueryMerchantCourierService,
		).
		Where("courier_id = ?", courierID).
		Find(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *merchantCourierServiceRepository) GetByMerchantIDCourierServiceID(merchantID uint, courierServiceID uint) (model *models.MerchantCourierService, err error) {
	result := r.db.
		Where("merchant_id = ?", merchantID).
		Where("courier_service_id = ?", courierServiceID).
		First(&model)
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

func (r *merchantCourierServiceRepository) UpsertUsingMap(modelMap map[string]any) (err error) {
	result := r.db.
		Model(&models.MerchantCourierService{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "courier_service_id"}, {Name: "merchant_id"}},
			DoUpdates: clause.Assignments(modelMap),
		}).Create(modelMap)
	if result.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (w *merchantCourierServiceRepository) Create(newModel *models.MerchantCourierService) (model *models.MerchantCourierService, err error) {
	queryResult := w.db.Create(&newModel)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if queryResult.RowsAffected == 0 {
		err = e.ErrInsertFailed
		return
	}

	model = newModel

	return
}

func (r *merchantCourierServiceRepository) GetAllCouriersByMerchantIDForProductUpdate(merchantID uint) (model []*internalProductDtos.Courier, err error) {
	result := r.db.
		Table("merchant_courier_services mcs").
		Select(`
			c.id AS id,
			c.code AS code,
			c.logo_url AS logo_url,
			c.name AS name
		`).
		Joins("JOIN courier_services cs ON  mcs.courier_service_id = cs.id ").
		Joins("JOIN couriers c ON cs.courier_id = c.id ").
		Where("merchant_id = ?", merchantID).
		Group("c.id").
		Find(&model)
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

func (r *merchantCourierServiceRepository) GetAllCourierServicesForProductUpdate(courierID uint, merchantID uint, productID uint) (model []*internalProductDtos.CourierService, err error) {
	result := r.db.
		Table("merchant_courier_services mcs").
		Select(`
			cs.id AS id,
			cs.code AS code,
			cs.description AS description,
			(
				CASE
					WHEN pcs.id IS NOT NULL
						THEN TRUE
					ELSE FALSE
				END
			) AS is_active
		`).
		Joins("JOIN courier_services cs ON  mcs.courier_service_id = cs.id ").
		Joins("LEFT JOIN product_courier_services pcs ON mcs.id = pcs.merchant_courier_service_id AND (?) = pcs.product_id", productID).
		Where("cs.courier_id = ?", courierID).
		Where("merchant_id = ?", merchantID).
		Find(&model)
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

func (r *merchantCourierServiceRepository) GetByConditionMap(conditionMap map[string]any) (model *models.MerchantCourierService, err error) {
	result := r.db.Where(conditionMap).First(&model)
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
