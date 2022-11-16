package repositories

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MerchantVoucherRepository interface {
	Create(newModel *models.MerchantVoucher) (model *models.MerchantVoucher, err error)
	GetDetailsByID(id uint) (*dtos.MerchantVoucherDetails, error)
	Update(newModel *models.MerchantVoucher) (model *models.MerchantVoucher, err error)
	Delete(id uint) (model *models.MerchantVoucher, err error)
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
		Where("end_datetime::date >= ?::date", time.Now()).
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

func (r *merchantVoucherRepository) GetDetailsByID(id uint) (*dtos.MerchantVoucherDetails, error) {
	var merchantVoucherDetails *dtos.MerchantVoucherDetails

	queryResult := r.db.Model(&models.MerchantVoucher{}).
		Select(`
			id, 
			merchant_id, 
			code, 
			start_datetime, 
			end_datetime, 
			discount_type, 
			discount, 
			minimum_spend_amount, 
			maximum_discount_amount, 
			maximum_quota, 
			current_quota, 
			(
				CASE WHEN NOW()::timestamp BETWEEN start_datetime::timestamp 
				AND end_datetime::timestamp THEN 'ongoing' 
				WHEN start_datetime::timestamp > NOW()::timestamp  THEN 'upcoming' 
				WHEN end_datetime::timestamp < NOW()::timestamp  THEN 'ended' END
			) AS voucher_status
		`).
		Where("id = ?", id).
		First(&merchantVoucherDetails)

	if queryResult.RowsAffected == 0 {
		return nil, e.ErrVoucherNotFound
	}

	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	return merchantVoucherDetails, nil
}

func (r *merchantVoucherRepository) Update(newModel *models.MerchantVoucher) (model *models.MerchantVoucher, err error) {
	queryResult := r.db.Clauses(clause.Returning{}).Updates(newModel)

	if queryResult.RowsAffected == 0 {
		err = e.ErrVoucherNotFound
		return
	}

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel
	return
}

func (r *merchantVoucherRepository) Delete(id uint) (model *models.MerchantVoucher, err error) {
	queryResult := r.db.Clauses(clause.Returning{}).Delete(&model, id)

	if queryResult.RowsAffected == 0 {
		err = e.ErrVoucherNotFound
		return
	}

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}
