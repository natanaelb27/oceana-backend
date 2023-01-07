package repositories

import (
	"github.com/natanaelb27/oceana-backend/internal/models"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	"gorm.io/gorm"
)

type DistrictRepository interface {
	GetAll(query *dtos.SubDistrictQuery) (model []*models.District, err error)
}

type districtRepository struct {
	db *gorm.DB
}

type DistrictRepositoryConfig struct {
	DB *gorm.DB
}

func NewDistrictRepository(c *DistrictRepositoryConfig) DistrictRepository {
	return &districtRepository{
		db: c.DB,
	}
}

func (e *districtRepository) GetAll(query *dtos.SubDistrictQuery) (model []*models.District, err error) {
	result := e.db.Model(&models.District{}).
		Joins("LEFT JOIN cities ON districts.city_id = cities.id").
		Where("districts.city_id = ?", query.City).
		Find(&model)

	err = result.Error
	return
}
