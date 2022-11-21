package repositories

import (
	"fmt"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MerchantVoucherRepository interface {
	Create(newModel *models.MerchantVoucher) (model *models.MerchantVoucher, err error)
	GetAllDetails(query *externalDtos.MerchantVoucherQuery, pagination *utils.Pagination) ([]*dtos.MerchantVoucherDetails, *utils.Pagination, error)
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

func (r *merchantVoucherRepository) GetAllDetails(query *externalDtos.MerchantVoucherQuery, pagination *utils.Pagination) ([]*dtos.MerchantVoucherDetails, *utils.Pagination, error) {
	var merchantVoucherDetails []*dtos.MerchantVoucherDetails

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
		`)

	if query.SearchVoucherCode != "" {
		queryResult.Where("code = ?", query.SearchVoucherCode)
	}

	if query.FilterVoucherStatus == "ongoing" {
		queryResult.Where("NOW()::timestamp BETWEEN start_datetime::timestamp AND end_datetime::timestamp")
	}

	if query.FilterVoucherStatus == "upcoming" {
		queryResult.Where("start_datetime::timestamp > NOW()::timestamp")
	}

	if query.FilterVoucherStatus == "ended" {
		queryResult.Where("end_datetime::timestamp < NOW()::timestamp")
	}

	if query.MerchantID != 0 {
		queryResult.Where("merchant_id = ?", query.MerchantID)
	}

	queryResult.
		Order(fmt.Sprintf("%s %s", query.SortBy, query.Sort)).
		Scopes(utils.Paginate(pagination, queryResult)).
		Find(&merchantVoucherDetails)

	if queryResult.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return merchantVoucherDetails, pagination, nil
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

	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrVoucherNotFound
		}
		return nil, e.ErrDatabaseOperationFailed
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
