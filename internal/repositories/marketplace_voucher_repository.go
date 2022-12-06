package repositories

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MarketplaceVoucherRepository interface {
	Create(newModel *models.MarketplaceVoucher) (model *models.MarketplaceVoucher, err error)
	GetDetailsByID(id uint) (*dtos.MarketplaceVoucherDetails, error)
	Update(newModel *models.MarketplaceVoucher) (model *models.MarketplaceVoucher, err error)
	Delete(id uint) (model *models.MarketplaceVoucher, err error)
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
	queryResult := r.db.
		Where("end_datetime::timestamp >= ?::timestamptz", time.Now()).
		FirstOrCreate(&newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel
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

func (r *marketplaceVoucherRepository) Update(newModel *models.MarketplaceVoucher) (model *models.MarketplaceVoucher, err error) {
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

func (r *marketplaceVoucherRepository) Delete(id uint) (model *models.MarketplaceVoucher, err error) {
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
