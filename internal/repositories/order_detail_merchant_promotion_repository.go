package repositories

import (
	"log"

	"github.com/natanaelb27/oceana-backend/internal/models"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

type OrderDetailMerchantPromotionRepository interface {
	WithTx(tx *gorm.DB) *orderDetailMerchantPromotionRepository
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

func (u orderDetailMerchantPromotionRepository) WithTx(tx *gorm.DB) *orderDetailMerchantPromotionRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
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
