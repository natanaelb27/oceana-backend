package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type MerchantPromotionRepository interface {
	CreateMyShopPromotion(merchantPromotion *models.MerchantPromotion) (*models.MerchantPromotion, error)
	GetAll(merchantID uint, query *dtos.MerchantPromotionQuery) (models []*models.MerchantPromotion, err error)
}

type merchantPromotionRepository struct {
	db *gorm.DB
}

type MerchantPromotionRepositoryConfig struct {
	DB *gorm.DB
}

func NewMerchantPromotionRepository(c *MerchantPromotionRepositoryConfig) MerchantPromotionRepository {
	return &merchantPromotionRepository{
		db: c.DB,
	}
}

func (mp *merchantPromotionRepository) CreateMyShopPromotion(newModel *models.MerchantPromotion) (model *models.MerchantPromotion, err error) {
	queryResult := mp.db.Create(&newModel)
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

func (mp *merchantPromotionRepository) GetAll(merchantID uint, query *dtos.MerchantPromotionQuery) (model []*models.MerchantPromotion, err error) {
	queryResult := mp.db.Model(&models.MerchantPromotion{}).Select(`
		id,
		merchant_id,
		name,
		start_datetime,
		end_datetime,
		(
			CASE WHEN NOW()::timestamp BETWEEN start_datetime::timestamp 
			AND end_datetime::timestamp THEN 'ongoing'
			WHEN start_datetime::timestamp > NOW()::timestamp  THEN 'upcoming' 
			WHEN end_datetime::timestamp < NOW()::timestamp  THEN 'ended' END
		) AS promotion_status
	`)

	if query.SearchPromotionName != "" {
		queryResult.Where("name = ?", query.SearchPromotionName)
	}

	if query.FilterPromotionStatus == "ongoing" {
		queryResult.Where("NOW()::timestamp BETWEEN start_datetime::timestamp AND end_datetime::timestamp")
	}

	if query.FilterPromotionStatus == "upcoming" {
		queryResult.Where("start_datetime::timestamp > NOW()::timestamp")
	}

	if query.FilterPromotionStatus == "ended" {
		queryResult.Where("end_datetime::timestamp < NOW()::timestamp")
	}

	queryResult.Where("merchant_id = ?", merchantID).Find(&model)
	if queryResult.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	return

}
