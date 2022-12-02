package repositories

import (
	"database/sql"

	internalLogisticDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/logistic_dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type MerchantCourierServiceRepository interface {
	GetAllByCourierID(courierID uint, merchantID uint) (model []*internalLogisticDtos.CourierService, err error)
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
	queryResult := r.db.Debug().
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
