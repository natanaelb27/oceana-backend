package repositories

import (
	"database/sql"
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetAllCategories() (model []*models.Category, err error)
	GetAllSubCategoriesByParentID(ID uint) (model []*models.Category, err error)
	GetAllByID(ID uint) (model []*models.Category, err error)
	GetByID(ID uint) (model *models.Category, err error)
	GetAncestorsByID(ID uint) (model []*models.Category, err error)
	WithTx(tx *gorm.DB) *categoryRepository
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

func (u categoryRepository) WithTx(tx *gorm.DB) *categoryRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *categoryRepository) GetAllCategories() (model []*models.Category, err error) {
	queryResult := r.db.Order("name ASC").Where("parent_id IS NULL").Find(&model)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *categoryRepository) GetAllSubCategoriesByParentID(ID uint) (model []*models.Category, err error) {
	queryResult := r.db.Order("name ASC").Where("parent_id = ?", ID).Find(&model)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *categoryRepository) GetAllByID(ID uint) (model []*models.Category, err error) {
	queryResult := r.db.Order("name ASC").Find(&model, ID)
	if queryResult.Error == gorm.ErrRecordNotFound {
		err = e.ErrDataNotFound
		return
	}
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *categoryRepository) GetByID(ID uint) (model *models.Category, err error) {
	queryResult := r.db.First(&model, ID)
	if queryResult.Error == gorm.ErrRecordNotFound {
		err = e.ErrDataNotFound
		return
	}
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *categoryRepository) GetAncestorsByID(ID uint) (model []*models.Category, err error) {
	queryResult := r.db.Raw(`
		WITH RECURSIVE recursive_category AS (
			SELECT 
				id,
				media_url, 
				name,
				parent_id 
				FROM 
				categories 
			WHERE 
				id = @category_id 
			UNION ALL 
			SELECT 
				categories.id,
				categories.media_url, 
				categories.name, 
				categories.parent_id 
			FROM 
				categories 
				JOIN recursive_category ON categories.id = recursive_category.parent_id
		) 
		SELECT 
			id, 
			media_url,
			name,
			parent_id
		FROM 
			recursive_category
		ORDER BY ROW_NUMBER () OVER () DESC
	`, sql.Named("category_id", ID)).
		Scan(&model)

	if queryResult.Error == gorm.ErrRecordNotFound {
		err = e.ErrDataNotFound
		return
	}
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}
