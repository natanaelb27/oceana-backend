package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SealabsPayRepository interface {
	Create(sealabsPay *models.SealabsPay) (model *models.SealabsPay, err error)
	GetAll(userID uuid.UUID) (model []*models.SealabsPay, err error)
	GetById(id uuid.UUID) (model *models.SealabsPay, err error)
}

type sealabsPayRepository struct {
	db *gorm.DB
}

type SealabsPayRepositoryConfig struct {
	DB *gorm.DB
}

func NewSealabsPayRepository(c *SealabsPayRepositoryConfig) SealabsPayRepository {
	return &sealabsPayRepository{
		db: c.DB,
	}
}

func (r *sealabsPayRepository) Create(newModel *models.SealabsPay) (model *models.SealabsPay, err error) {
	queryResult := r.db.Where("account_number = ?", newModel.AccountNumber).FirstOrCreate(&newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	if queryResult.RowsAffected == 0 {
		err = e.ErrSealabsPayAccountNumberAlreadyExist
		return
	}

	model = newModel

	return
}

func (r *sealabsPayRepository) GetAll(userID uuid.UUID) (model []*models.SealabsPay, err error) {
	queryResult := r.db.Where("user_id = ?", userID).Find(&model)
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *sealabsPayRepository) GetById(id uuid.UUID) (model *models.SealabsPay, err error) {
	queryResult := r.db.First(&model, id)
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}
