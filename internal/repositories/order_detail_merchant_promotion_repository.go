package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type OrderDetailMerchantPromotionRepository interface {
	Create(newModel []*models.OrderDetailMerchantPromotion) (model []*models.OrderDetailMerchantPromotion, err error)
}

type orderDetailMerchantPromotionRepository struct {
	db *gorm.DB
}

type OrderDetailMerchantPromotionRepositoryConfig struct {
	DB *gorm.DB
}

func NewOrderDetailMerchantPromotionRepository(c *OrderDetailMerchantPromotionRepositoryConfig) OrderDetailMerchantPromotionRepository {
	return &orderDetailMerchantPromotionRepository{
		db: c.DB,
	}
}

func (r *orderDetailMerchantPromotionRepository) Create(newModel []*models.OrderDetailMerchantPromotion) (model []*models.OrderDetailMerchantPromotion, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}
