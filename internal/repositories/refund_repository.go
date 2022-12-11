package repositories

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type RefundRepository interface {
	UpdateAllByStatusByRefundExpirationDatetimeLteNow(status string, datas *models.Refund) (err error)
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

func (r *refundRepository) UpdateAllByStatusByRefundExpirationDatetimeLteNow(status string, datas *models.Refund) (err error) {
	queryResult := r.db.
		Where("status = ?", status).
		Where("refund_expiration_datetime <= ?", time.Now()).
		Updates(&datas)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}
