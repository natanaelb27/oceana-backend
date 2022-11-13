package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *models.Role) (model *models.Role, err error)
	GetAll() (model []*models.Role, err error)
	GetById(id uuid.UUID) (model *models.Role, err error)
	Update(id uuid.UUID, newModel *models.Role) (model *models.Role, err error)
	Delete(id uuid.UUID) (model *models.Role, err error)
}

type roleRepository struct {
	db *gorm.DB
}

type RoleRepositoryConfig struct {
	DB *gorm.DB
}

func NewRoleRepository(c *RoleRepositoryConfig) RoleRepository {
	return &roleRepository{
		db: c.DB,
	}
}

func (r *roleRepository) Create(newModel *models.Role) (model *models.Role, err error) {
	queryResult := r.db.Create(&newModel)
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

func (r *roleRepository) GetAll() (model []*models.Role, err error) {
	queryResult := r.db.Find(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *roleRepository) GetById(id uuid.UUID) (model *models.Role, err error) {
	result := r.db.First(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *roleRepository) Update(id uuid.UUID, newModel *models.Role) (model *models.Role, err error) {
	result := r.db.Model(&model).Where("id = ?", id).Updates(&newModel)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	model.ID = id

	return
}

func (r *roleRepository) Delete(id uuid.UUID) (model *models.Role, err error) {
	result := r.db.Delete(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model.ID = id

	return
}
