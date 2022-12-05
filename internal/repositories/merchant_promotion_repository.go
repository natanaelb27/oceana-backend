package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type MerchantPromotionRepository interface {
	CreateMyShopPromotion(merchantPromotion *models.MerchantPromotion, productIDList []int64) (*models.MerchantPromotion, error)
	GetAll(merchantID uint, query *dtos.MerchantPromotionQuery) (models []*models.MerchantPromotion, err error)
	GetByID(id uint) (*models.MerchantPromotion, error)
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

func (mp *merchantPromotionRepository) CreateMyShopPromotion(newModel *models.MerchantPromotion, productIDList []int64) (model *models.MerchantPromotion, err error) {
	subQueryProductList := mp.db.Model(&models.ProductMerchantPromotion{}).
		Select("ARRAY_AGG(product_merchant_promotions.product_id)").
		Joins("JOIN merchant_promotions ON merchant_promotions.id = product_merchant_promotions.merchant_promotion_id").
		Where("merchant_promotions.start_datetime::timestamp <= ?::timestamptz", newModel.EndDatetime).
		Where("?::timestamptz <= merchant_promotions.end_datetime::timestamp", newModel.StartDatetime).
		Where("merchant_promotions.merchant_id = ?", newModel.MerchantId)

	queryResult := mp.db.
		Where("(?) && (?)", pq.Int64Array(productIDList), subQueryProductList).
		FirstOrCreate(&newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	if queryResult.RowsAffected == 0 {
		err = e.ErrMerchantPromotionProductConflict
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

func (mp *merchantPromotionRepository) GetByID(id uint) (*models.MerchantPromotion, error) {
	var merchantPromotion *models.MerchantPromotion

	result := mp.db.First(&merchantPromotion, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrPromotionNotFound
		}
		return nil, e.ErrDatabaseOperationFailed
	}

	return merchantPromotion, nil
}
