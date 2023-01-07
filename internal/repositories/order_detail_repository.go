package repositories

import (
	"log"

	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/constants/product_media_constants"
	"github.com/natanaelb27/oceana-backend/internal/dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

type OrderDetailRepository interface {
	WithTx(tx *gorm.DB) *orderDetailRepository
	OrderDetailListGetByOrderID(id uuid.UUID) ([]*dtos.OrderDetailList, error)
	Create(newModel []*models.OrderDetail) (model []*models.OrderDetail, err error)
	GetByID(id uuid.UUID) (model *models.OrderDetail, err error)
	GetDetailsByID(id uuid.UUID) (model *models.OrderDetail, err error)
}

type orderDetailRepository struct {
	db *gorm.DB
}

type OrderDetailRepositoryConfig struct {
	DB *gorm.DB
}

func NewOrderDetailRepository(c *OrderDetailRepositoryConfig) OrderDetailRepository {
	return &orderDetailRepository{
		db: c.DB,
	}
}

func (u orderDetailRepository) WithTx(tx *gorm.DB) *orderDetailRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *orderDetailRepository) OrderDetailListGetByOrderID(id uuid.UUID) ([]*dtos.OrderDetailList, error) {
	var orderDetailList []*dtos.OrderDetailList

	subQueryVariants := r.db.Table("product_snapshots AS ps").
		Select(`
			STRING_AGG(product_variation_snapshots.variant, ', ')
		`).
		Joins("JOIN product_variation_snapshots ON ps.id = product_variation_snapshots.product_snapshot_id").
		Where("ps.id = product_snapshots.id").
		Group("ps.id")

	subQueryProductImage := r.db.Table("product_snapshots AS ps").
		Select(`
			product_media_snapshots.media_url
		`).
		Joins("JOIN product_media_snapshots ON ps.id = product_media_snapshots.product_snapshot_id").
		Where("ps.id = product_snapshots.id").
		Where("product_media_snapshots.type = ?", product_media_constants.TypeImage).
		Limit(1)

	subQueryReviewImage := r.db.Model(&models.ProductReviewMedia{}).
		Select(`
			product_reviews.order_detail_id,
			ARRAY_AGG(product_review_media.media_url) AS review_images
		`).
		Joins("JOIN product_reviews ON product_reviews.id = product_review_media.product_review_id").
		Where("type = ?", product_media_constants.TypeImage).
		Group("product_reviews.order_detail_id")

	subQueryReviewVideo := r.db.Model(&models.ProductReviewMedia{}).
		Select(`
			product_reviews.order_detail_id,
			ARRAY_AGG(product_review_media.media_url) AS review_videos
		`).
		Joins("JOIN product_reviews ON product_reviews.id = product_review_media.product_review_id").
		Where("type = ?", product_media_constants.TypeVideo).
		Group("product_reviews.order_detail_id")

	queryResult := r.db.Model(&models.OrderDetail{}).
		Select(`
			order_details.id AS order_detail_id,
			order_details.message AS message,
			product_snapshots.product_id AS product_id,
			(?) AS product_image,
			product_snapshots.name AS product_name,
			products.name AS current_product_name,
			(?) AS product_variations,
			order_details.quantity AS quantity,
			product_snapshots.price AS price,
			order_details.total_price AS total_price,
			product_reviews.id AS product_review_id,
			product_reviews.rating AS rating,
			product_reviews.review AS review,
			tbl_review_image.review_images AS review_images,
			tbl_review_video.review_videos AS review_videos,
			product_reviews.created_at AS review_created_datetime
		`, subQueryProductImage,
			subQueryVariants).
		Joins("JOIN orders ON order_details.order_id = orders.id").
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN product_snapshots ON order_details.product_snapshot_id = product_snapshots.id").
		Joins("JOIN products ON products.id = product_snapshots.product_id").
		Joins("LEFT JOIN product_reviews ON order_details.id = product_reviews.order_detail_id").
		Joins("LEFT JOIN (?) tbl_review_image ON order_details.id = tbl_review_image.order_detail_id", subQueryReviewImage).
		Joins("LEFT JOIN (?) tbl_review_video ON order_details.id = tbl_review_video.order_detail_id", subQueryReviewVideo).
		Where("orders.id = ?", id).
		Find(&orderDetailList)

	if queryResult.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	return orderDetailList, nil
}

func (r *orderDetailRepository) Create(newModel []*models.OrderDetail) (model []*models.OrderDetail, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}

func (r *orderDetailRepository) GetByID(id uuid.UUID) (model *models.OrderDetail, err error) {
	queryResult := r.db.First(&model, id)

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

func (r *orderDetailRepository) GetDetailsByID(id uuid.UUID) (model *models.OrderDetail, err error) {
	queryResult := r.db.Model(&models.OrderDetail{}).
		Preload("ProductSnapshot").
		Preload("Order").
		Preload("Order.Transaction").
		First(&model, id)

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
