package repositories

import (
	"fmt"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MarketplaceVoucherRepository interface {
	Create(marketplaceVoucher *models.MarketplaceVoucher) (model *models.MarketplaceVoucher, err error)
	GetAll() (model []*models.MarketplaceVoucher, err error)
	GetAllByActive() (model []*models.MarketplaceVoucher, err error)
	GetById(id uint) (model *models.MarketplaceVoucher, err error)
	Update(id uint, newModel *models.MarketplaceVoucher) (err error)
	Delete(id uint) (model *models.MarketplaceVoucher, err error)
	GetDetailsByID(id uint) (*dtos.MarketplaceVoucherDetails, error)
	AddCurrentQuota(ID uint, amount uint) (*models.MarketplaceVoucher, error)
	GetAllDetailsByBuyer(id uuid.UUID, query *externalDtos.BuyerMarketplaceVoucherQuery, pagination *utils.Pagination) ([]*dtos.MarketplaceVoucherDetails, *utils.Pagination, error)
}

type marketplaceVoucherRepository struct {
	db *gorm.DB
}

type MarketplaceVoucherRepositoryConfig struct {
	DB *gorm.DB
}

func NewMarketplaceVoucherRepository(c *MarketplaceVoucherRepositoryConfig) MarketplaceVoucherRepository {
	return &marketplaceVoucherRepository{
		db: c.DB,
	}
}

func (r *marketplaceVoucherRepository) Create(newModel *models.MarketplaceVoucher) (model *models.MarketplaceVoucher, err error) {
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

func (r *marketplaceVoucherRepository) GetAllByActive() (model []*models.MarketplaceVoucher, err error) {
	queryResult := r.db.Debug().
		Where("NOW()::timestamp BETWEEN start_datetime::timestamp AND end_datetime::timestamp").
		Find(&model)
	if queryResult.Error != nil && queryResult.Error != e.ErrDataNotFound {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if queryResult.Error == e.ErrDataNotFound {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *marketplaceVoucherRepository) GetAll() (model []*models.MarketplaceVoucher, err error) {
	queryResult := r.db.Find(&model)
	if queryResult.Error != nil && queryResult.Error != e.ErrDataNotFound {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if queryResult.Error == e.ErrDataNotFound {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *marketplaceVoucherRepository) GetById(id uint) (model *models.MarketplaceVoucher, err error) {
	queryResult := r.db.First(&model, id)
	if queryResult.Error != nil && queryResult.Error != e.ErrDataNotFound {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if queryResult.Error == e.ErrDataNotFound {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *marketplaceVoucherRepository) Update(id uint, newModel *models.MarketplaceVoucher) (err error) {
	queryResult := r.db.Model(&models.MarketplaceVoucher{}).Where("id = ?", id).Updates(&newModel)
	if queryResult.Error != nil && queryResult.Error != e.ErrDataNotFound {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if queryResult.Error == e.ErrDataNotFound {
		err = e.ErrDataNotFound
		return
	}
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *marketplaceVoucherRepository) Delete(id uint) (model *models.MarketplaceVoucher, err error) {
	queryResult := r.db.Delete(&model, id)
	if queryResult.Error != nil && queryResult.Error != e.ErrDataNotFound {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if queryResult.Error == e.ErrDataNotFound {
		err = e.ErrDataNotFound
		return
	}
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model.ID = id

	return
}

func (r *marketplaceVoucherRepository) GetDetailsByID(id uint) (*dtos.MarketplaceVoucherDetails, error) {
	var marketplaceVoucherDetails *dtos.MarketplaceVoucherDetails

	queryResult := r.db.Model(&models.MarketplaceVoucher{}).
		Select(`
			id, 
			name, 
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
		First(&marketplaceVoucherDetails)

	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrVoucherNotFound
		}
		return nil, e.ErrDatabaseOperationFailed
	}

	return marketplaceVoucherDetails, nil
}

func (r *marketplaceVoucherRepository) AddCurrentQuota(ID uint, amount uint) (*models.MarketplaceVoucher, error) {
	var model *models.MarketplaceVoucher

	result := r.db.
		Model(&models.MarketplaceVoucher{}).
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

func (r *marketplaceVoucherRepository) GetAllDetailsByBuyer(id uuid.UUID, query *externalDtos.BuyerMarketplaceVoucherQuery, pagination *utils.Pagination) ([]*dtos.MarketplaceVoucherDetails, *utils.Pagination, error) {
	var marketplaceVoucherDetails []*dtos.MarketplaceVoucherDetails

	subQueryUnavailableVoucher := r.db.Model(&models.Order{}).
		Select(`
			order_merchant_vouchers.merchant_voucher_id
		`).
		Joins("JOIN transactions ON transactions.id = orders.transaction_id").
		Joins("JOIN order_merchant_vouchers ON order_merchant_vouchers.order_id = orders.id").
		Where("transactions.user_id = ?", id)

	queryResult := r.db.Model(&models.MarketplaceVoucher{}).
		Select(`
			id, 
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

	if query.FilterVoucherAvailability == "available" {
		queryResult.Where("merchant_vouchers.id NOT IN (?)", subQueryUnavailableVoucher)
	}

	if query.FilterVoucherAvailability == "unavailable" {
		queryResult.Where("merchant_vouchers.id IN (?)", subQueryUnavailableVoucher)
	}

	queryResult.
		Order(fmt.Sprintf("%s %s", query.SortBy, query.Sort)).
		Scopes(utils.Paginate(pagination, queryResult)).
		Find(&marketplaceVoucherDetails)

	if queryResult.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return marketplaceVoucherDetails, pagination, nil
}
