package repositories

import (
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type CourierRepository interface {
	WithTx(tx *gorm.DB) *courierRepository
	GetAll() (model []*models.Courier, err error)
	GetByCode(code string) (model *models.Courier, err error)
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

func (u courierRepository) WithTx(tx *gorm.DB) *courierRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *courierRepository) GetAll() (model []*models.Courier, err error) {
	queryResult := r.db.
		Order("id ASC").
		Find(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *courierRepository) GetByCode(code string) (model *models.Courier, err error) {
	queryResult := r.db.Where("code = ?", code).First(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}
	return
}
