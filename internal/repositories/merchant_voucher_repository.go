package repositories

import (
	"fmt"
	"log"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MerchantVoucherRepository interface {
	WithTx(tx *gorm.DB) *merchantVoucherRepository
	Create(newModel *models.MerchantVoucher) (model *models.MerchantVoucher, err error)
	GetAllDetails(query *externalDtos.MerchantVoucherQuery, pagination *utils.Pagination) ([]*dtos.MerchantVoucherDetails, *utils.Pagination, error)
	GetDetailsByID(id uint) (*dtos.MerchantVoucherDetails, error)
	Update(newModel *models.MerchantVoucher) (model *models.MerchantVoucher, err error)
	Delete(id uint) (model *models.MerchantVoucher, err error)
	ValidateMerchantVoucherByBuyer(merchantVoucherID uint, userID uuid.UUID) (resultMerchantVoucherID *uint, err error)
	GetAllDetailsByBuyer(id uuid.UUID, query *externalDtos.BuyerMerchantVoucherQuery, pagination *utils.Pagination) ([]*dtos.MerchantVoucherDetails, *utils.Pagination, error)
	AddCurrentQuota(ID uint, amount uint) (*models.MerchantVoucher, error)
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

func (u merchantVoucherRepository) WithTx(tx *gorm.DB) *merchantVoucherRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *merchantVoucherRepository) Create(newModel *models.MerchantVoucher) (model *models.MerchantVoucher, err error) {
	queryResult := r.db.
		Where("code = ?", newModel.Code).
		Where("end_datetime::timestamp >= ?::timestamptz", time.Now()).
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

	if query.FilterVoucherStatus == constants.OngoingPromotionOrVoucher {
		queryResult.Where("NOW()::timestamp BETWEEN start_datetime::timestamp AND end_datetime::timestamp")
	}

	if query.FilterVoucherStatus == constants.UpcomingPromotionOrVoucher {
		queryResult.Where("start_datetime::timestamp > NOW()::timestamp")
	}

	if query.FilterVoucherStatus == constants.EndedPromotionOrVoucher {
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

func (r *merchantVoucherRepository) GetAllDetailsByBuyer(id uuid.UUID, query *externalDtos.BuyerMerchantVoucherQuery, pagination *utils.Pagination) ([]*dtos.MerchantVoucherDetails, *utils.Pagination, error) {
	var merchantVoucherDetails []*dtos.MerchantVoucherDetails

	subQueryUnavailableVoucher := r.db.Model(&models.Order{}).
		Select(`
			order_merchant_vouchers.merchant_voucher_id
		`).
		Joins("JOIN transactions ON transactions.id = orders.transaction_id").
		Joins("JOIN order_merchant_vouchers ON order_merchant_vouchers.order_id = orders.id").
		Where("transactions.user_id = ?", id)

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

	if query.FilterVoucherAvailability == "available" {
		queryResult.Where("merchant_vouchers.id NOT IN (?)", subQueryUnavailableVoucher)
	}

	if query.FilterVoucherAvailability == "unavailable" {
		queryResult.Where("merchant_vouchers.id IN (?)", subQueryUnavailableVoucher)
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

func (r *merchantVoucherRepository) ValidateMerchantVoucherByBuyer(merchantVoucherID uint, userID uuid.UUID) (resultMerchantVoucherID *uint, err error) {
	queryResult := r.db.Debug().Model(&models.Order{}).
		Select(`
			order_merchant_vouchers.merchant_voucher_id
		`).
		Joins("JOIN transactions ON transactions.id = orders.transaction_id").
		Joins("JOIN order_merchant_vouchers ON order_merchant_vouchers.order_id = orders.id").
		Where("transactions.user_id = ?", userID).
		Where("order_merchant_vouchers.merchant_voucher_id = ?", merchantVoucherID).
		Scan(&resultMerchantVoucherID)

	if queryResult.Error != nil {
		if queryResult.RowsAffected == 0 {
			err = e.ErrDataNotFound
			return
		}

		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *merchantVoucherRepository) AddCurrentQuota(ID uint, amount uint) (*models.MerchantVoucher, error) {
	var model *models.MerchantVoucher

	result := r.db.
		Model(&models.MerchantVoucher{}).
		Where("id = ?", ID).
		Update("current_quota", gorm.Expr("current_quota + ?", amount))

	if result.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrDataNotFound
	}

	return model, nil
}
