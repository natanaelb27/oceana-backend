package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type RefundRepository interface {
	Create(newModel *models.Refund) (model *models.Refund, err error)
}

type refundRepository struct {
	db *gorm.DB
}

type RefundRepositoryConfig struct {
	DB *gorm.DB
}

func NewRefundRepository(c *RefundRepositoryConfig) RefundRepository {
	return &refundRepository{
		db: c.DB,
	}
}

func (r *refundRepository) Create(newModel *models.Refund) (model *models.Refund, err error) {
	result := r.db.
		Where("order_id = ?", newModel.OrderID).
		Where("status IN (?, ?, ?, ?)", constants.RefundStatusOngoing, constants.RefundStatusApproved, constants.RefundStatusApprovedBySystem, constants.RefundStatusCompleted).
		FirstOrCreate(&newModel)

	if result.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	if result.RowsAffected == 0 {
		err = e.ErrRefundAlreadyExist
		return
	}

	model = newModel

	return
}
