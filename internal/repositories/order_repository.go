package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	OrderListGetAll(id uuid.UUID, query *externalDtos.OrderListQuery, pagination *utils.Pagination) ([]*dtos.OrderList, *utils.Pagination, error)
	OrderListDetailsGetAll(id uuid.UUID) (*dtos.OrderListDetails, error)
}

type orderRepository struct {
	db *gorm.DB
}

type OrderRepositoryConfig struct {
	DB *gorm.DB
}

func NewOrderRepository(c *OrderRepositoryConfig) OrderRepository {
	return &orderRepository{
		db: c.DB,
	}
}

func (r *orderRepository) OrderListGetAll(id uuid.UUID, query *externalDtos.OrderListQuery, pagination *utils.Pagination) ([]*dtos.OrderList, *utils.Pagination, error) {
	var orderList []*dtos.OrderList

	queryResult := r.db.Model(&models.Order{}).
		Select(`
			orders.transaction_id AS transaction_id,
			orders.id AS order_id,
			orders.merchant_id AS merchant_id,
			merchants.name AS shop_name,
			transactions.transaction_date AS transaction_datetime,
			orders.status AS order_status,
			(orders.total_price_after_discount::float / transactions.total_price::float * transactions.total_price_after_discount) AS total_price
		`).
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Where("transactions.user_id = ?", id)

	if query.FilterOrderStatus == "ongoing" {
		if query.FilterOngoingOrderStatus != "" && query.FilterOngoingOrderStatus != "all" {
			queryResult.Where("orders.status = ?", query.FilterOngoingOrderStatus)
		} else {
			queryResult.Where("orders.status IN ('pending payment', 'pending seller', 'out for delivery', 'delivered', 'complained')")
		}
	}

	if query.FilterOrderStatus == "completed" {
		queryResult.Where("orders.status = 'completed'")
	}

	if query.FilterOrderStatus == "canceled" {
		queryResult.Where("orders.status = 'canceled'")
	}

	if query.FilterStartTransactionDate != nil && query.FilterEndTransactionDate != nil {
		queryResult.Where("transactions.transaction_date::date BETWEEN ?::date AND ?::date", query.FilterStartTransactionDate, query.FilterEndTransactionDate)
	}

	queryResult.Scopes(utils.Paginate(pagination, queryResult)).
		Find(&orderList)

	if queryResult.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return orderList, pagination, nil
}

func (r *orderRepository) OrderListDetailsGetAll(id uuid.UUID) (*dtos.OrderListDetails, error) {
	var orderListDetails *dtos.OrderListDetails

	subQueryTotalPriceBeforeDiscount := r.db.Model(&models.OrderDetail{}).
		Select(`
				SUM(order_details.total_price)
			`).
		Where("order_details.order_id = orders.id").
		Group("order_details.order_id")

	subQueryMerchantPromotion := r.db.Model(&models.OrderDetailMerchantPromotion{}).
		Select(`
			SUM(order_detail_merchant_promotions.total_discount) AS total_merchant_promotion_discount
	`).
		Joins("JOIN order_details ON order_detail_merchant_promotions.order_detail_id = order_details.id").
		Where("order_details.order_id = orders.id").
		Group("order_details.order_id")

	subQueryMerchantVoucher := r.db.Model(&models.OrderMerchantVoucher{}).
		Select(`
			SUM(order_merchant_vouchers.total_discount) AS total_merchant_voucher_discount
	`).
		Where("order_merchant_vouchers.order_id = orders.id").
		Group("order_merchant_vouchers.order_id")

	subQueryMarketplacePromotion := r.db.Model(&models.TransactionMarketplacePromotion{}).
		Select(`
			SUM(transaction_marketplace_promotions.total_discount)
		`).
		Where("transaction_marketplace_promotions.transaction_id = transactions.id").
		Group("transaction_marketplace_promotions.transaction_id")

	subQueryMarketplaceVoucher := r.db.Model(&models.TransactionMarketplaceVoucher{}).
		Select(`
				SUM(transaction_marketplace_vouchers.total_discount)
			`).
		Where("transaction_marketplace_vouchers.transaction_id = transactions.id").
		Group("transaction_marketplace_vouchers.transaction_id")

	queryResult := r.db.Model(&models.Order{}).
		Select(`
			orders.transaction_id AS transaction_id,
			transactions.user_id AS user_id,
			orders.id AS order_id,
			orders.merchant_id AS merchant_id,
			merchants.name AS shop_name,
			transactions.transaction_date AS transaction_datetime,
			orders.status AS order_status,
			courier_services.description AS courier_description,
			payment_methods.name AS payment_method,	
			(?) AS total_price_before_discount,
			(?) AS total_merchant_promotion_discount,
			(?) AS total_merchant_voucher_discount,
			orders.total_price_after_discount::float / transactions.total_price::float * (?) AS total_marketplace_promotion_discount,
			orders.total_price_after_discount::float / transactions.total_price::float * (?) AS total_marketplace_voucher_discount,
			orders.shipping_fee AS shipping_fee,
			(orders.total_price_after_discount::float / transactions.total_price::float * transactions.total_price_after_discount) + orders.shipping_fee AS total_payment
		`,
			subQueryTotalPriceBeforeDiscount,
			subQueryMerchantPromotion,
			subQueryMerchantVoucher,
			subQueryMarketplacePromotion,
			subQueryMarketplaceVoucher).
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN payment_methods ON transactions.payment_method_id = payment_methods.id").
		Joins("JOIN courier_services ON orders.courier_service_id = courier_services.id").
		Where("orders.id = ?", id).
		First(&orderListDetails)

	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrOrderNotFound
		}
		return nil, e.ErrDatabaseOperationFailed
	}

	return orderListDetails, nil
}
