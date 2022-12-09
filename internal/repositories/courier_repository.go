package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type CourierRepository interface {
	GetAll() (model []*models.Courier, err error)
}

type courierRepository struct {
	db *gorm.DB
}

type CourierRepositoryConfig struct {
	DB *gorm.DB
}

func NewCourierRepository(c *CourierRepositoryConfig) CourierRepository {
	return &courierRepository{
		db: c.DB,
	}
}

func (r *courierRepository) GetAll() (model []*models.Courier, err error) {
	queryResult := r.db.Find(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}
