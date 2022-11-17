package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type SealabsPayRepository interface {
	Create(sealabsPay *models.SealabsPay) (model *models.SealabsPay, err error)
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
