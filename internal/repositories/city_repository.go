package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"gorm.io/gorm"
)

type CityRepository interface {
	GetAll(query *dtos.SubDistrictQuery) (model []*models.City, err error)
}

type cityRepository struct {
	db *gorm.DB
}

type CityRepositoryConfig struct {
	DB *gorm.DB
}

func NewCityRepository(c *CityRepositoryConfig) CityRepository {
	return &cityRepository{
		db: c.DB,
	}
}

func (e *cityRepository) GetAll(query *dtos.SubDistrictQuery) (model []*models.City, err error) {
	result := e.db.Model(&models.City{}).
		Joins("LEFT JOIN provinces ON cities.province_id = provinces.id").
		Where("cities.province_id = ?", query.Province).
		Find(&model)

	err = result.Error
	return
}
