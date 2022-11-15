package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

type UserRepositoryConfig struct {
	DB *gorm.DB
}

func NewUserRepository(c *UserRepositoryConfig) UserRepository {
	return &userRepository{
		db: c.DB,
	}
}

func (u *userRepository) Create(newModel *models.User) (model *models.User, err error) {
	queryResult := u.db.Create(&newModel)
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
