package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AddressRepository interface {
	GetProfileAddressess(userID uuid.UUID) (model []*models.Address, err error)
	Create(address *models.Address) (*models.Address, error)
	GetProfileAddressById(id uuid.UUID) (*models.Address, error)
	UpdateProfileAddress(id uuid.UUID, newModel *models.Address) (model *models.Address, err error)
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

func (a *addressRepository) GetProfileAddressess(userID uuid.UUID) (model []*models.Address, err error) {
	queryResult := a.db.Joins("JOIN users on users.id = addresses.user_id").Joins("JOIN sub_districts on sub_districts.id = addresses.sub_district_id").Joins("JOIN districts on districts.id = sub_districts.district_id").Joins("JOIN cities on cities.id = districts.city_id").Joins("JOIN provinces on provinces.id = cities.province_id").Preload("User").Preload("SubDistrict").Preload("SubDistrict.District").Preload("SubDistrict.District.City").Preload("SubDistrict.District.City.Province").Where("user_id = ?", userID).Find(&model)

	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (a *addressRepository) GetProfileAddressById(id uuid.UUID) (model *models.Address, err error) {
	queryResult := a.db.Joins("JOIN users on users.id = addresses.user_id").Joins("JOIN sub_districts on sub_districts.id = addresses.sub_district_id").Joins("JOIN districts on districts.id = sub_districts.district_id").Joins("JOIN cities on cities.id = districts.city_id").Joins("JOIN provinces on provinces.id = cities.province_id").Preload("User").Preload("SubDistrict").Preload("SubDistrict.District").Preload("SubDistrict.District.City").Preload("SubDistrict.District.City.Province").First(&model, id)

	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (a *addressRepository) UpdateProfileAddress(id uuid.UUID, newModel *models.Address) (model *models.Address, err error) {
	result := a.db.Model(&model).Where("id = ?", id).Updates(&newModel)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	model.ID = id

	return
}
