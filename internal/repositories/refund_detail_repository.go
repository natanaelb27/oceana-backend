package repositories

import (
	"database/sql"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefundDetailRepository interface {
	CalculateRefundDetails(orderDetailID uuid.UUID, quantity uint) (model *models.RefundDetail, err error)
	Create(newModel []*models.RefundDetail) (model []*models.RefundDetail, err error)
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
