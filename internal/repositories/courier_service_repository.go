package repositories

import (
	"log"

	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type CourierServiceRepository interface {
	WithTx(tx *gorm.DB) *courierServiceRepository
	GetAllByCourierID(courierID uint) (model []*internalDtos.CourierService, err error)
	GetByCourierIDAndServiceCode(courierID uint, courierServiceCode string) (model *internalDtos.CourierService, err error)
}

type courierServiceRepository struct {
	db *gorm.DB
}

type CourierServiceRepositoryConfig struct {
	DB *gorm.DB
}

func NewCourierServiceRepository(c *CourierServiceRepositoryConfig) CourierServiceRepository {
	return &courierServiceRepository{
		db: c.DB,
	}
}

func (u courierServiceRepository) WithTx(tx *gorm.DB) *courierServiceRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *courierServiceRepository) GetAllByCourierID(courierID uint) (model []*internalDtos.CourierService, err error) {
	queryResult := r.db.Where("courier_id = ?", courierID).Find(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *courierServiceRepository) GetByCourierIDAndServiceCode(courierID uint, courierServiceCode string) (model *internalDtos.CourierService, err error) {
	queryResult := r.db.
		Where("courier_id = ? AND code = ?", courierID, courierServiceCode).
		First(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}
	return
}
