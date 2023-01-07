package repositories

import (
	"log"

	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/models"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserMerchantPromotionRepository interface {
	WithTx(tx *gorm.DB) *userMerchantPromotionRepository
	GetUserUsage(userID uuid.UUID, productMerchantPromotionID uint) (model *models.UserMerchantPromotion, err error)
	UpsertUsingMap(modelMap map[string]any) (model *models.UserMerchantPromotion, err error)
}

type userMerchantPromotionRepository struct {
	db *gorm.DB
}

type UserMerchantPromotionRepositoryConfig struct {
	DB *gorm.DB
}

func NewUserMerchantPromotionRepository(c *UserMerchantPromotionRepositoryConfig) UserMerchantPromotionRepository {
	return &userMerchantPromotionRepository{
		db: c.DB,
	}
}

func (u userMerchantPromotionRepository) WithTx(tx *gorm.DB) *userMerchantPromotionRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *userMerchantPromotionRepository) GetUserUsage(userID uuid.UUID, productMerchantPromotionID uint) (model *models.UserMerchantPromotion, err error) {
	queryResult := r.db.
		Where("user_id = ?", userID).
		Where("product_merchant_promotion_id = ?", productMerchantPromotionID).
		First(&model)

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

func (r *userMerchantPromotionRepository) UpsertUsingMap(modelMap map[string]any) (model *models.UserMerchantPromotion, err error) {
	queryResult := r.db.
		Model(&models.UserMerchantPromotion{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "product_merchant_promotion_id"}},
			DoUpdates: clause.Assignments(modelMap),
		}).Create(modelMap)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}
