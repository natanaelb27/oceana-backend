package repositories

import (
	"database/sql"
	"strings"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
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
	OrderListDetailsGetByID(id uuid.UUID) (*dtos.OrderListDetails, error)
	GetByID(id uuid.UUID) (model *models.Order, err error)
	Update(id uuid.UUID, newModel *models.Order) (model *models.Order, err error)
	GetOrderCount() (*int64, error)
	Create(newModel []*models.Order) (model []*models.Order, err error)
	OrderListGetAllByMerchantID(merchantID uint, query *externalDtos.OrderListQuery, pagination *utils.Pagination) ([]*dtos.OrderList, *utils.Pagination, error)
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
			orders.invoice_number AS invoice_number,
			orders.id AS order_id,
			orders.merchant_id AS merchant_id,
			merchants.name AS shop_name,
			users.id AS user_id,
			users.full_name AS user_full_name,
			courier_services.id AS courier_service_id,
			courier_services.description AS courier_description,
			transactions.transaction_datetime AS transaction_datetime,
			orders.status AS order_status,
			(orders.total_price_after_discount::float / transactions.total_price::float * transactions.total_price_after_discount) + orders.shipping_fee AS total_price,
			refunds.id AS refund_id

		`).
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN users ON users.id = transactions.user_id").
		Joins("JOIN courier_services ON orders.courier_service_id = courier_services.id").
		Joins("LEFT JOIN refunds ON orders.id = refunds.order_id AND refunds.status = ?", constants.RefundStatusCompleted).
		Where("transactions.user_id = ?", id).
		Where("LOWER(orders.status) != LOWER(?)", constants.PendingPaymentStatus)

	if query.Search != "" {
		subQuerySearch := r.db.Model(&models.Order{}).
			Select("DISTINCT orders.id").
			Joins("JOIN merchants ON merchants.id = orders.merchant_id").
			Joins("JOIN order_details ON order_details.order_id = orders.id").
			Joins("JOIN product_snapshots ON product_snapshots.id = order_details.product_snapshot_id").
			Where("merchants.name ILIKE @search OR product_snapshots.name ILIKE @search", sql.Named("search", "%"+query.Search+"%"))

		queryResult.Where("orders.id IN (?)", subQuerySearch)
	}

	if strings.EqualFold(query.FilterOrderStatus, constants.OngoingStatus) {
		if query.FilterOngoingOrderStatus != "" && !strings.EqualFold(query.FilterOngoingOrderStatus, constants.AllStatus) {
			queryResult.Where("orders.status ILIKE ?", query.FilterOngoingOrderStatus)
		} else {
			queryResult.Where("LOWER(orders.status) IN (LOWER(?), LOWER(?), LOWER(?), LOWER(?), LOWER(?))",
				constants.PendingSellerStatus,
				constants.OnProcessStatus,
				constants.OutForDeliveryStatus,
				constants.DeliveredStatus,
				constants.ComplainedStatus)
		}
	} else if strings.EqualFold(query.FilterOrderStatus, constants.CanceledStatus) {
		queryResult.Where("orders.status ILIKE ?", "%canceled%")
	} else {
		if query.FilterOrderStatus != "" && !strings.EqualFold(query.FilterOrderStatus, constants.AllStatus) {
			queryResult.Where("orders.status ILIKE ?", query.FilterOrderStatus)
		}
	}

	if query.FilterStartTransactionDate != nil && query.FilterEndTransactionDate != nil {
		queryResult.Where("transactions.transaction_datetime::date BETWEEN ?::date AND ?::date", query.FilterStartTransactionDate, query.FilterEndTransactionDate)
	}

	queryResult.Scopes(utils.Paginate(pagination, queryResult)).
		Find(&orderList)

	if queryResult.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return orderList, pagination, nil
}

func (r *orderRepository) OrderListDetailsGetByID(id uuid.UUID) (*dtos.OrderListDetails, error) {
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

	subQueryMarketplaceVoucher := r.db.Model(&models.TransactionMarketplaceVoucher{}).
		Select(`
				SUM(transaction_marketplace_vouchers.total_discount)
			`).
		Where("transaction_marketplace_vouchers.transaction_id = transactions.id").
		Group("transaction_marketplace_vouchers.transaction_id")

	queryResult := r.db.Model(&models.Order{}).
		Select(`
			orders.transaction_id AS transaction_id,
			users.id AS user_id,
			users.full_name AS user_full_name,
			orders.invoice_number AS invoice_number,
			orders.id AS order_id,
			orders.merchant_id AS merchant_id,
			merchants.name AS shop_name,
			transactions.transaction_datetime AS transaction_datetime,
			orders.status AS order_status,
			courier_services.id AS courier_service_id,
			courier_services.description AS courier_description,
			payment_methods.name AS payment_method,	
			(?) AS total_price_before_discount,
			(?) AS total_merchant_promotion_discount,
			(?) AS total_merchant_voucher_discount,
			orders.total_price_after_discount::float / transactions.total_price::float * (?) AS total_marketplace_voucher_discount,
			orders.shipping_fee AS shipping_fee,
			(orders.total_price_after_discount::float / transactions.total_price::float * transactions.total_price_after_discount) + orders.shipping_fee AS total_payment,
			refunds.id AS refund_id
		`,
			subQueryTotalPriceBeforeDiscount,
			subQueryMerchantPromotion,
			subQueryMerchantVoucher,
			subQueryMarketplaceVoucher).
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN users ON users.id = transactions.user_id").
		Joins("JOIN payment_methods ON transactions.payment_method_id = payment_methods.id").
		Joins("JOIN courier_services ON orders.courier_service_id = courier_services.id").
		Joins("LEFT JOIN refunds ON orders.id = refunds.order_id AND refunds.status = ?", constants.RefundStatusCompleted).
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

func (r *orderRepository) GetByID(id uuid.UUID) (model *models.Order, err error) {
	queryResult := r.db.First(&model, id)
	if queryResult.Error == gorm.ErrRecordNotFound {
		err = e.ErrDataNotFound
		return
	}
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *orderRepository) Update(id uuid.UUID, newModel *models.Order) (model *models.Order, err error) {
	queryResult := r.db.Model(&model).Where("id = ?", id).Updates(&newModel)
	if queryResult.Error == gorm.ErrRecordNotFound {
		err = e.ErrDataNotFound
		return
	}
	if queryResult.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	model.ID = id

	return
}

func (r *orderRepository) GetOrderCount() (*int64, error) {
	var count int64
	queryResult := r.db.Model(&models.Order{}).Count(&count)

	if queryResult.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	return &count, nil
}

func (r *orderRepository) Create(newModel []*models.Order) (model []*models.Order, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}

func (r *orderRepository) OrderListGetAllByMerchantID(merchantID uint, query *externalDtos.OrderListQuery, pagination *utils.Pagination) ([]*dtos.OrderList, *utils.Pagination, error) {
	var orderList []*dtos.OrderList

	queryResult := r.db.Model(&models.Order{}).
		Select(`
			orders.transaction_id AS transaction_id,
			orders.invoice_number AS invoice_number,
			orders.id AS order_id,
			orders.merchant_id AS merchant_id,
			merchants.name AS shop_name,
			users.id AS user_id,
			users.full_name AS user_full_name,
			courier_services.id AS courier_service_id,
			courier_services.description AS courier_description,
			transactions.transaction_datetime AS transaction_datetime,
			orders.status AS order_status,
			orders.total_price_after_discount::float AS total_price,
			refunds.id AS refund_id
		`).
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN users ON users.id = transactions.user_id").
		Joins("JOIN courier_services ON orders.courier_service_id = courier_services.id").
		Joins("LEFT JOIN refunds ON orders.id = refunds.order_id AND refunds.status = ?", constants.RefundStatusCompleted).
		Where("orders.merchant_id = ?", merchantID)
		// Where("LOWER(orders.status) != LOWER(?)", constants.PendingPaymentStatus)

	if query.Search != "" {
		subQuerySearch := r.db.Model(&models.Order{}).
			Select("DISTINCT orders.id").
			Joins("JOIN merchants ON merchants.id = orders.merchant_id").
			Joins("JOIN order_details ON order_details.order_id = orders.id").
			Joins("JOIN product_snapshots ON product_snapshots.id = order_details.product_snapshot_id").
			Where("merchants.name ILIKE @search OR product_snapshots.name ILIKE @search", sql.Named("search", "%"+query.Search+"%"))

		queryResult.Where("orders.id IN (?)", subQuerySearch)
	}

	if strings.EqualFold(query.FilterOrderStatus, constants.OngoingStatus) {
		if query.FilterOngoingOrderStatus != "" && !strings.EqualFold(query.FilterOngoingOrderStatus, constants.AllStatus) {
			queryResult.Where("orders.status ILIKE ?", query.FilterOngoingOrderStatus)
		} else {
			queryResult.Where("LOWER(orders.status) IN (LOWER(?), LOWER(?), LOWER(?), LOWER(?), LOWER(?), LOWER(?))",
				constants.PendingPaymentStatus,
				constants.PendingSellerStatus,
				constants.OnProcessStatus,
				constants.OutForDeliveryStatus,
				constants.DeliveredStatus,
				constants.ComplainedStatus)
		}
	} else if strings.EqualFold(query.FilterOrderStatus, constants.CanceledStatus) {
		queryResult.Where("orders.status ILIKE ?", "%canceled%")
	} else {
		if query.FilterOrderStatus != "" && !strings.EqualFold(query.FilterOrderStatus, constants.AllStatus) {
			queryResult.Where("orders.status ILIKE ?", query.FilterOrderStatus)
		}
	}

	if query.FilterStartTransactionDate != nil && query.FilterEndTransactionDate != nil {
		queryResult.Where("transactions.transaction_datetime::date BETWEEN ?::date AND ?::date", query.FilterStartTransactionDate, query.FilterEndTransactionDate)
	}

	queryResult.Scopes(utils.Paginate(pagination, queryResult)).
		Find(&orderList)

	if queryResult.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return orderList, pagination, nil
}
