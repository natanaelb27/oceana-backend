package repositories

import (
	"database/sql"

	internalLogisticDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/logistic_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MerchantCourierServiceRepository interface {
	GetAllByCourierID(courierID uint, merchantID uint) (model []*internalLogisticDtos.CourierService, err error)
	GetByMerchantIDCourierServiceID(merchantID uint, courierServiceID uint) (model *models.MerchantCourierService, err error)
	UpsertUsingMap(modelMap map[string]any) (err error)
	Create(newModel *models.MerchantCourierService) (model *models.MerchantCourierService, err error)
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
