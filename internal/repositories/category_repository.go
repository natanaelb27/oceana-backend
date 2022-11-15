package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetAllCategories() (model []*models.Category, err error)
	GetAllSubCategoriesByParentID(ID uint) (model []*models.Category, err error)
}

type categoryRepository struct {
	db *gorm.DB
}

type CategoryRepositoryConfig struct {
	DB *gorm.DB
}

func NewCategoryRepository(c *CategoryRepositoryConfig) CategoryRepository {
	return &categoryRepository{
		db: c.DB,
	}
}

func (r *categoryRepository) GetAllCategories() (model []*models.Category, err error) {
	queryResult := r.db.Where("parent_id IS NULL").Find(&model)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *categoryRepository) GetAllSubCategoriesByParentID(ID uint) (model []*models.Category, err error) {
	queryResult := r.db.Where("parent_id = ?", ID).Find(&model)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}
