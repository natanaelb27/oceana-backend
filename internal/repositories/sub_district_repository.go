package repositories

import (
	"github.com/natanaelb27/oceana-backend/internal/models"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	"gorm.io/gorm"
)

type SubDistrictRepository interface {
	GetAll(query *dtos.SubDistrictQuery) (model []*models.SubDistrict, err error)
}

type subDistrictRepository struct {
	db *gorm.DB
}

type SubDistrictRepositoryConfig struct {
	DB *gorm.DB
}

func NewSubDistrictRepository(c *SubDistrictRepositoryConfig) SubDistrictRepository {
	return &subDistrictRepository{
		db: c.DB,
	}
}

func (e *subDistrictRepository) GetAll(query *dtos.SubDistrictQuery) (model []*models.SubDistrict, err error) {
	result := e.db.Model(&models.SubDistrict{}).
		Joins("LEFT JOIN districts ON sub_districts.district_id = districts.id").
		Where("sub_districts.district_id = ?", query.District).
		Find(&model)

	err = result.Error
	return
}
