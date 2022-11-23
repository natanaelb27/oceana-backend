package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	Create(merchant *models.Merchant) (*models.Merchant, error)
	GetByName(name string) (*models.Merchant, error)
	GetByUserID(id uuid.UUID) (model *models.Merchant, err error)
	GetInformationByID(id uint) (merchantInformation *dtos.MerchantInformation, err error)
}

type merchantRepository struct {
	db *gorm.DB
}

type MerchantRepositoryConfig struct {
	DB *gorm.DB
}

func NewMerchantRepository(c *MerchantRepositoryConfig) MerchantRepository {
	return &merchantRepository{
		db: c.DB,
	}
}

func (m *merchantRepository) Create(newModel *models.Merchant) (model *models.Merchant, err error) {
	queryResult := m.db.Create(&newModel)
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

func (m *merchantRepository) GetByName(name string) (model *models.Merchant, err error) {
	result := m.db.Model(&models.Merchant{}).Where("name = ?", name).Find(&model)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}
	return
}

func (r *merchantRepository) GetByUserID(id uuid.UUID) (model *models.Merchant, err error) {
	result := r.db.Where("user_id = ?", id).First(&model)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			err = e.ErrMerchantNotFound
			return
		}

		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *merchantRepository) GetInformationByID(id uint) (merchantInformation *dtos.MerchantInformation, err error) {
	subQueryProduct := r.db.Model(&models.Merchant{}).
		Select(`
			merchants.id AS merchant_id,
			COUNT(products.id) AS product_count
		`).
		Joins("JOIN products ON merchants.id = products.merchant_id").
		Where("merchants.id = ?", id).
		Group("merchants.id")

	// subQueryMerchantReview := r.db.Model(&models.Merchant{}).
	// 	Select(`
	// 		merchants.id AS merchant_id,
	// 		COUNT(product_reviews.id) AS total_rating,
	// 		AVG(product_reviews.rating) AS rating_star
	// 	`).
	// 	Joins("JOIN orders ON merchants.id = orders.merchant_id").
	// 	Joins("JOIN order_details ON orders.id = order_details.order_id").
	// 	Joins("JOIN product_reviews ON order_details.id = product_reviews.order_detail_id").
	// 	Where("merchants.id = ?", id).
	// 	Group("merchants.id")

	subQueryMerchantReview := r.db.Model(&models.ProductReview{}).
		Select(`
			products.merchant_id AS merchant_id,
			COUNT(product_reviews.id) AS total_rating,
			AVG(product_reviews.rating) AS rating_star
		`).
		Joins("JOIN products ON product_reviews.product_id = products.id").
		Where("products.merchant_id = ?", id).
		Group("products.merchant_id")

	result := r.db.Model(&models.Merchant{}).
		Select(`
			merchants.id AS id,
			merchants.name AS name,
			merchants.photo AS photo,
			merchant_review_agg.total_rating AS total_rating,
			merchant_review_agg.rating_star AS rating_star,
			product_agg.product_count AS product_count,
			merchants.user_id,
			addresses.detail AS shop_address,
			sub_districts.name AS shop_sub_district,
			districts.name AS shop_district,
			cities.name AS shop_city,
			provinces.name AS shop_province
		`).
		Joins("LEFT JOIN (?) product_agg ON merchants.id = product_agg.merchant_id", subQueryProduct).
		Joins("LEFT JOIN (?) merchant_review_agg ON merchants.id = merchant_review_agg.merchant_id", subQueryMerchantReview).
		Joins("JOIN addresses ON merchants.address_id = addresses.id").
		Joins("JOIN sub_districts ON addresses.sub_district_id = sub_districts.id").
		Joins("JOIN districts ON sub_districts.district_id = districts.id").
		Joins("JOIN cities ON districts.city_id = cities.id").
		Joins("JOIN provinces ON cities.province_id = provinces.id").
		Where("merchants.id = ?", id).
		First(&merchantInformation)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			err = e.ErrMerchantNotFound
			return
		}

		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}
