package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type RefundMediaRepository interface {
	Create(newModel []*models.RefundMedia) (model []*models.RefundMedia, err error)
}

type refundMediaRepository struct {
	db *gorm.DB
}

type RefundMediaRepositoryConfig struct {
	DB *gorm.DB
}

func NewRefundMediaRepository(c *RefundMediaRepositoryConfig) RefundMediaRepository {
	return &refundMediaRepository{
		db: c.DB,
	}
}

func (r *refundMediaRepository) Create(newModel []*models.RefundMedia) (model []*models.RefundMedia, err error) {
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
