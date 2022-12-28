package repositories

import (
	"database/sql"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants/product_media_constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefundDetailRepository interface {
	CalculateRefundDetails(orderDetailID uuid.UUID, quantity uint) (model *models.RefundDetail, err error)
	Create(newModel []*models.RefundDetail) (model []*models.RefundDetail, err error)
	GetDetailsByRefundID(refundID uuid.UUID) ([]*dtos.RefundItems, error)
}

type refundDetailRepository struct {
	db *gorm.DB
}

type RefundDetailRepositoryConfig struct {
	DB *gorm.DB
}

func NewRefundDetailRepository(c *RefundDetailRepositoryConfig) RefundDetailRepository {
	return &refundDetailRepository{
		db: c.DB,
	}
}

func (r *refundDetailRepository) CalculateRefundDetails(orderDetailID uuid.UUID, quantity uint) (model *models.RefundDetail, err error) {
	subQueryMerchantLost := r.db.Model(&models.OrderDetail{}).
		Select(`
			order_details.id AS order_detail_id,
			@refund_quantity * (order_details.total_price_after_discount / order_details.quantity) / orders.total_price * orders.total_price_after_discount AS amount_merchant_lost`,
			sql.Named("refund_quantity", quantity)).
		Joins("JOIN orders ON order_details.order_id = orders.id").
		Where("order_details.id = @order_detail_id", sql.Named("order_detail_id", orderDetailID))

	subQueryUserGet := r.db.Model(&models.OrderDetail{}).
		Select(`
			order_details.id AS order_detail_id,
			(@refund_quantity * (order_details.total_price_after_discount / orders.total_price * orders.total_price_after_discount) / order_details.quantity) / transactions.total_price * transactions.total_price_after_discount AS amount_user_get`,
			sql.Named("refund_quantity", quantity)).
		Joins("JOIN orders ON order_details.order_id = orders.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Where("order_details.id = @order_detail_id", sql.Named("order_detail_id", orderDetailID))

	result := r.db.Model(&models.OrderDetail{}).
		Select(`
			order_details.id AS order_detail_id,
			user_get.amount_user_get AS amount_user_get,
			merchant_lost.amount_merchant_lost AS amount_merchant_lost
		`).
		Joins("JOIN (?) merchant_lost ON order_details.id = merchant_lost.order_detail_id", subQueryMerchantLost).
		Joins("JOIN (?) user_get ON order_details.id = user_get.order_detail_id", subQueryUserGet).
		Where("order_details.id = @order_detail_id", sql.Named("order_detail_id", orderDetailID)).
		First(&model)

	if result.Error != nil {
		if result.RowsAffected == 0 {
			err = e.ErrDataNotFound
			return
		}

		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *refundDetailRepository) Create(newModel []*models.RefundDetail) (model []*models.RefundDetail, err error) {
	result := r.db.Create(&newModel)

	if result.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}

func (r *refundDetailRepository) GetDetailsByRefundID(refundID uuid.UUID) ([]*dtos.RefundItems, error) {
	var refundItems []*dtos.RefundItems

	subQueryVariants := r.db.Model(&models.ProductVariationSnapshot{}).
		Select(`
			product_variation_snapshots.product_snapshot_id AS product_snapshot_id,
			STRING_AGG(product_variation_snapshots.variant, ', ') AS product_variations
		`).
		Group("product_variation_snapshots.product_snapshot_id")

	subQueryProductImage := r.db.Table("product_snapshots AS ps").
		Select(`
			product_media_snapshots.media_url
		`).
		Joins("JOIN product_media_snapshots ON ps.id = product_media_snapshots.product_snapshot_id").
		Where("ps.id = product_snapshots.id").
		Where("product_media_snapshots.type = ?", product_media_constants.TypeImage).
		Limit(1)

	result := r.db.Model(&models.RefundDetail{}).
		Select(`
			refund_details.id AS refund_detail_id,
			refund_details.order_detail_id AS order_detail_id,
			(?) AS product_image,
			product_snapshots.name AS product_name,
			tbl_variation.product_variations AS product_variations,
			refund_details.quantity AS refund_quantity,
			refund_details.amount_user_get AS user_refund_amount,
			refund_details.amount_merchant_lost AS merchant_lost_amount
		`, subQueryProductImage).
		Joins("JOIN order_details ON refund_details.order_detail_id = order_details.id").
		Joins("JOIN product_snapshots ON order_details.product_snapshot_id = product_snapshots.id").
		Joins("LEFT JOIN (?) tbl_variation ON product_snapshots.id = tbl_variation.product_snapshot_id", subQueryVariants).
		Where("refund_details.refund_id = ?", refundID).
		Find(&refundItems)

	if result.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	return refundItems, nil
}
