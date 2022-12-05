package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderDetailRepository interface {
	OrderDetailListGetByOrderID(id uuid.UUID) ([]*dtos.OrderDetailList, error)
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
		Limit(1)

	queryResult := r.db.Model(&models.OrderDetail{}).
		Select(`
			order_details.id AS order_detail_id,
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
