package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AddressRepository interface {
	GetProfileAddressess(userID uuid.UUID) (addressDetails []*dtos.AddressDetails, err error)
	Create(address *models.Address) (*models.Address, error)
	GetProfileAddressById(id uuid.UUID) (*dtos.AddressDetails, error)
	UpdateProfileAddress(id uuid.UUID, newModel *models.Address) (model *models.Address, err error)
	DeleteProfileAddress(id uuid.UUID) (model *models.Address, err error)
	UpdateProfileAddressDefault(id uuid.UUID, userID uuid.UUID) (model *models.Address, err error)
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
func (a *addressRepository) GetProfileAddressess(userID uuid.UUID) (addressDetails []*dtos.AddressDetails, err error) {
	queryResult := a.db.Model(&models.Address{}).
		Select(`
			addresses.id AS id,
			addresses.user_id AS user_id,
			users.phone_number AS phone_number,
			addresses.name AS name,
			addresses.detail AS detail,
			sub_districts.id AS sub_district_id,
			sub_districts.name AS sub_district,
			sub_districts.zip_code AS zip_code,
			districts.id AS district_id,
			districts.name AS district,
			cities.id AS city_id,
			cities.name AS city,
			provinces.id AS province_id,
			provinces.name AS province,
			addresses.is_default AS is_default,
			(CASE 
				WHEN merchants.id IS NOT NULL THEN TRUE 
				ELSE FALSE 
			END) AS is_shop_address
		`).
		Joins("JOIN users on addresses.user_id = users.id").
		Joins("LEFT JOIN merchants on merchants.address_id = addresses.id").
		Joins("JOIN sub_districts on sub_districts.id = addresses.sub_district_id").
		Joins("JOIN districts on districts.id = sub_districts.district_id").
		Joins("JOIN cities on cities.id = districts.city_id").
		Joins("JOIN provinces on provinces.id = cities.province_id").
		Where("addresses.user_id = ?", userID).
		Find(&addressDetails)

	if queryResult.Error != nil {
		err = queryResult.Error
		return
	}

	return
}

func (a *addressRepository) GetProfileAddressById(id uuid.UUID) (addressDetails *dtos.AddressDetails, err error) {
	queryResult := a.db.Model(&models.Address{}).
		Select(`
			addresses.id AS id,
			addresses.user_id AS user_id,
			users.phone_number AS phone_number,
			addresses.name AS name,
			addresses.detail AS detail,
			sub_districts.id AS sub_district_id,
			sub_districts.name AS sub_district,
			sub_districts.zip_code AS zip_code,
			districts.id AS district_id,
			districts.name AS district,
			cities.id AS city_id,
			cities.name AS city,
			provinces.id AS province_id,
			provinces.name AS province,
			addresses.is_default AS is_default,
			(CASE 
				WHEN merchants.id IS NOT NULL THEN TRUE 
				ELSE FALSE 
			END) AS is_shop_address
		`).
		Joins("JOIN users on addresses.user_id = users.id").
		Joins("LEFT JOIN merchants on merchants.address_id = addresses.id").
		Joins("JOIN sub_districts on sub_districts.id = addresses.sub_district_id").
		Joins("JOIN districts on districts.id = sub_districts.district_id").
		Joins("JOIN cities on cities.id = districts.city_id").
		Joins("JOIN provinces on provinces.id = cities.province_id").
		Where("addresses.id = ?", id).
		First(&addressDetails)

	if queryResult.Error != nil {
		err = queryResult.Error
		return
	}

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

func (a *addressRepository) DeleteProfileAddress(id uuid.UUID) (model *models.Address, err error) {
	result := a.db.Delete(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model.ID = id

	return
}

func (a *addressRepository) UpdateProfileAddressDefault(id uuid.UUID, userID uuid.UUID) (model *models.Address, err error) {
	updateExisting := a.db.Model(&model).Where("user_id = ?", userID).Where("is_default = ?", true).Update("is_default", false)
	if updateExisting.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	result := a.db.Model(&model).Where("id = ?", id).Update("is_default", true)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	getUpdated := a.db.First(&model, id)
	if getUpdated.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}
