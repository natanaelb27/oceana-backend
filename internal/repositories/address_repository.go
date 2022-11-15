package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(address *models.Address) (*models.Address, error)
}

type addressRepository struct {
	db *gorm.DB
}

type AddressRepositoryConfig struct {
	DB *gorm.DB
}

func NewAddressRepository(c *AddressRepositoryConfig) AddressRepository {
	return &addressRepository{
		db: c.DB,
	}
}

func (a *addressRepository) Create(newModel *models.Address) (model *models.Address, err error) {
	queryResult := a.db.Create(&newModel)
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
