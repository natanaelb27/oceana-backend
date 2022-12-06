package repositories

import (
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type CourierServiceRepository interface {
	GetAllByCourierID(courierID uint) (model []*internalDtos.CourierService, err error)
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

func (r *courierServiceRepository) GetAllByCourierID(courierID uint) (model []*internalDtos.CourierService, err error) {
	queryResult := r.db.Where("courier_id = ?", courierID).Find(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}
