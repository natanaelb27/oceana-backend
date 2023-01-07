package repositories

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/constants"
	"github.com/natanaelb27/oceana-backend/internal/dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	externalDtos "github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

type MarketplaceVoucherRepository interface {
	WithTx(tx *gorm.DB) *marketplaceVoucherRepository
	Create(marketplaceVoucher *models.MarketplaceVoucher) (model *models.MarketplaceVoucher, err error)
	GetAll() (model []*models.MarketplaceVoucher, err error)
	GetAllDetails(query *externalDtos.MarketplaceVoucherQuery, pagination *utils.Pagination) ([]*models.MarketplaceVoucher, *utils.Pagination, error)
	GetAllByActive() (model []*models.MarketplaceVoucher, err error)
	GetById(id uint) (model *models.MarketplaceVoucher, err error)
	Update(id uint, newModel *models.MarketplaceVoucher) (err error)
	Delete(id uint) (model *models.MarketplaceVoucher, err error)
	GetDetailsByID(id uint) (*dtos.MarketplaceVoucherDetails, error)
	AddCurrentQuota(ID uint, amount uint) (*models.MarketplaceVoucher, error)
	GetAllDetailsByBuyer(id uuid.UUID, query *externalDtos.BuyerMarketplaceVoucherQuery, pagination *utils.Pagination) ([]*dtos.MarketplaceVoucherDetails, *utils.Pagination, error)
	ValidateMarketplaceVoucherByBuyer(marketplaceVoucherID uint, userID uuid.UUID) (resultMarketplaceVoucherID *uint, err error)
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

func (u marketplaceVoucherRepository) WithTx(tx *gorm.DB) *marketplaceVoucherRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
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
	queryResult := r.db.
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

func (r *marketplaceVoucherRepository) GetAllDetails(query *externalDtos.MarketplaceVoucherQuery, pagination *utils.Pagination) ([]*models.MarketplaceVoucher, *utils.Pagination, error) {
	var marketplaceVoucherDetails []*models.MarketplaceVoucher

	queryResult := r.db.Model(&models.MarketplaceVoucher{}).
		Select(`
			id, 
			banner_url, 
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
		`)

	if query.SearchVoucherName != "" {
		queryResult.Where("name ILIKE ?", "%"+query.SearchVoucherName+"%")
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

	queryResult.
		Order(fmt.Sprintf("%s %s", query.SortBy, query.Sort)).
		Scopes(utils.Paginate(pagination, queryResult)).
		Find(&marketplaceVoucherDetails)

	if queryResult.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return marketplaceVoucherDetails, pagination, nil
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
			banner_url, 
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

	subQueryUnavailableVoucher := r.db.Model(&models.Transaction{}).
		Select(`
			transaction_marketplace_vouchers.marketplace_voucher_id
		`).
		Joins("JOIN transaction_marketplace_vouchers ON transaction_marketplace_vouchers.transaction_id = transactions.id").
		Where("transactions.user_id = ?", id)

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
		`)

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
		queryResult.Where("marketplace_vouchers.id NOT IN (?)", subQueryUnavailableVoucher)
	}

	if query.FilterVoucherAvailability == "unavailable" {
		queryResult.Where("marketplace_vouchers.id IN (?)", subQueryUnavailableVoucher)
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

func (r *marketplaceVoucherRepository) ValidateMarketplaceVoucherByBuyer(marketplaceVoucherID uint, userID uuid.UUID) (resultMarketplaceVoucherID *uint, err error) {
	queryResult := r.db.Model(&models.Transaction{}).
		Select(`
			transaction_marketplace_vouchers.marketplace_voucher_id
		`).
		Joins("JOIN transaction_marketplace_vouchers ON transaction_marketplace_vouchers.transaction_id = transactions.id").
		Where("transactions.user_id = ?", userID).
		Where("transaction_marketplace_vouchers.marketplace_voucher_id = ?", marketplaceVoucherID).
		Scan(&resultMarketplaceVoucherID)

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
