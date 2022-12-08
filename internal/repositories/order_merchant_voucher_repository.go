package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type OrderMerchantVoucherRepository interface {
	Create(newModel []*models.OrderMerchantVoucher) (model []*models.OrderMerchantVoucher, err error)
}

type orderMerchantVoucherRepository struct {
	db *gorm.DB
}

type OrderMerchantVoucherRepositoryConfig struct {
	DB *gorm.DB
}

func NewOrderMerchantVoucherRepository(c *OrderMerchantVoucherRepositoryConfig) OrderMerchantVoucherRepository {
	return &orderMerchantVoucherRepository{
		db: c.DB,
	}
}

func (r *orderMerchantVoucherRepository) Create(newModel []*models.OrderMerchantVoucher) (model []*models.OrderMerchantVoucher, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}
