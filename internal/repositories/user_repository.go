package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetByPhoneNumber(phoneNumber string) (model *models.User, err error)
	GetByUsername(username string) (model *models.User, err error)
	GetByEmailID(emailID uuid.UUID) (model *models.User, err error)
	GetByID(id uuid.UUID) (model *models.User, err error)
	Update(id uuid.UUID, newModel *models.User) (model *models.User, err error)
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

func (u *userRepository) GetByPhoneNumber(phoneNumber string) (model *models.User, err error) {
	result := u.db.Model(&models.User{}).Where("phone_number = ?", phoneNumber).Find(&model)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}
	return
}

func (u *userRepository) GetByUsername(username string) (model *models.User, err error) {
	result := u.db.Model(&models.User{}).Where("username = ?", username).Find(&model)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}
	return
}

func (u *userRepository) GetByEmailID(emailID uuid.UUID) (model *models.User, err error) {
	result := u.db.Model(&models.User{}).Where("email_id = ?", emailID).Find(&model)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}
	return
}

func (u *userRepository) GetByID(id uuid.UUID) (model *models.User, err error) {
	result := u.db.First(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (u *userRepository) Update(id uuid.UUID, newModel *models.User) (model *models.User, err error) {
	result := u.db.Model(&model).Where("id = ?", id).Updates(&newModel)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	model.ID = id

	return
}
