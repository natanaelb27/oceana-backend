package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type MerchantVoucherRepository interface {
	Create(newModel *models.MerchantVoucher) (model *models.MerchantVoucher, err error)
}

type merchantVoucherRepository struct {
	db *gorm.DB
}

type MerchantVoucherRepositoryConfig struct {
	DB *gorm.DB
}

func NewMerchantVoucherRepository(c *MerchantVoucherRepositoryConfig) MerchantVoucherRepository {
	return &merchantVoucherRepository{
		db: c.DB,
	}
}

func (r *merchantVoucherRepository) Create(newModel *models.MerchantVoucher) (model *models.MerchantVoucher, err error) {
	queryResult := r.db.
		Where("code = ?", newModel.Code).
		Where("end_datetime::date >= now()::date").
		Where("merchant_id = ?", newModel.MerchantID).
		FirstOrCreate(&newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	if queryResult.RowsAffected == 0 {
		err = e.ErrMerchantVoucherCodeAlreadyExist
		return
	}

	model = newModel
	return
}
