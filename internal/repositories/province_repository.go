package repositories

import (
	"github.com/natanaelb27/oceana-backend/internal/models"
	"gorm.io/gorm"
)

type ProvinceRepository interface {
	GetAll() (model []*models.Province, err error)
}

type provinceRepository struct {
	db *gorm.DB
}

type ProvinceRepositoryConfig struct {
	DB *gorm.DB
}

func NewProvinceRepository(c *ProvinceRepositoryConfig) ProvinceRepository {
	return &provinceRepository{
		db: c.DB,
	}
}

func (e *provinceRepository) GetAll() (model []*models.Province, err error) {
	result := e.db.Find(&model)

	err = result.Error
	return
}
