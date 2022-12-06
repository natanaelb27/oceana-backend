package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type MarketplaceVoucherRepository interface {
	Create(marketplaceVoucher *models.MarketplaceVoucher) (model *models.MarketplaceVoucher, err error)
	GetAll() (model []*models.MarketplaceVoucher, err error)
	GetById(id uint) (model *models.MarketplaceVoucher, err error)
	Update(id uint, newModel *models.MarketplaceVoucher) (err error)
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
