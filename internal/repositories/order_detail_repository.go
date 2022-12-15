package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/constants/product_media_constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderDetailRepository interface {
	OrderDetailListGetByOrderID(id uuid.UUID) ([]*dtos.OrderDetailList, error)
	Create(newModel []*models.OrderDetail) (model []*models.OrderDetail, err error)
	GetByID(id uuid.UUID) (model *models.OrderDetail, err error)
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

	queryResult := r.db.Model(&models.OrderDetail{}).
		Select(`
			order_details.id AS order_detail_id,
			order_details.message AS message,
			(?) AS product_image,
			product_snapshots.name AS product_name,
			(?) AS product_variations,
			order_details.quantity AS quantity,
			product_snapshots.price AS price,
			order_details.total_price AS total_price
		`, subQueryProductImage,
			subQueryVariants).
		Joins("JOIN orders ON order_details.order_id = orders.id").
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN product_snapshots ON order_details.product_snapshot_id = product_snapshots.id").
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
