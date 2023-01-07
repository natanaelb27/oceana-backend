package repositories

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/constants"
	"github.com/natanaelb27/oceana-backend/internal/dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	externalDtos "github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

type OrderRepository interface {
	WithTx(tx *gorm.DB) *orderRepository
	OrderListGetAll(id uuid.UUID, query *externalDtos.OrderListQuery, pagination *utils.Pagination) ([]*dtos.OrderList, *utils.Pagination, error)
	OrderListDetailsGetByID(id uuid.UUID) (*dtos.OrderListDetails, error)
	GetByID(id uuid.UUID) (model *models.Order, err error)
	GetByTransactionID(transactionID uuid.UUID) (models []*models.Order, err error)
	Update(id uuid.UUID, newModel *models.Order) (model *models.Order, err error)
	GetOrderCount() (*int64, error)
	Create(newModel []*models.Order) (model []*models.Order, err error)
	OrderListGetAllByMerchantID(merchantID uint, query *externalDtos.OrderListQuery, pagination *utils.Pagination) ([]*dtos.OrderList, *utils.Pagination, error)
	PendingPaymentGetAll() (pendingOrders []*models.Order, err error)
	GetMaxProcessedDaysByOrderID(orderID uuid.UUID) (maxProcessedDays uint, err error)
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

func (u orderRepository) WithTx(tx *gorm.DB) *orderRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *orderRepository) OrderListGetAll(id uuid.UUID, query *externalDtos.OrderListQuery, pagination *utils.Pagination) ([]*dtos.OrderList, *utils.Pagination, error) {
	var orderList []*dtos.OrderList

	sellerConfirmationExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_SELLER_CONFIRMATION_EXPIRATION_INTERVAL"))
	if err != nil {
		sellerConfirmationExpirationInterval = constants.CronSellerConfirmationExpirationInterval
	}

	autoCompleteDeliveredExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_AUTO_COMPLETE_DELIVERED_EXPIRATION_INTERVAL"))
	if err != nil {
		autoCompleteDeliveredExpirationInterval = constants.CronAutoCompleteDeliveredExpirationInterval
	}

	autoCompleteOrderRejectedRefundExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_AUTO_COMPLETE_ORDER_REJECTED_REFUND_EXPIRATION_INTERVAL"))
	if err != nil {
		autoCompleteOrderRejectedRefundExpirationInterval = constants.CronAutoCompleteOrderRejectedRefundExpirationInterval
	}

	autoApproveRefundSellerConfirmationExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_AUTO_APPROVE_REFUND_SELLER_CONFIRMATION_EXPIRATION_INTERVAL"))
	if err != nil {
		autoApproveRefundSellerConfirmationExpirationInterval = constants.CronAutoApproveRefundSellerConfirmationExpirationInterval
	}

	subQueryRefundRank := r.db.Model(&models.Refund{}).
		Select(`
			id AS refund_id,
			order_id AS order_id,
			status AS status,
			refund_request_datetime AS refund_request_datetime,
			rejected_datetime AS rejected_datetime,
			RANK() OVER (PARTITION BY order_id ORDER BY refund_request_datetime DESC) AS refund_rank
		`)

	subQueryLastRefundStatus := r.db.Table("(?) AS tbl_refund_rank", subQueryRefundRank).
		Select(`
			tbl_refund_rank.refund_id AS refund_id,
			tbl_refund_rank.order_id AS order_id,
			tbl_refund_rank.status AS status,
			tbl_refund_rank.refund_request_datetime AS refund_request_datetime,
			tbl_refund_rank.rejected_datetime AS rejected_datetime
		`).
		Where("tbl_refund_rank.refund_rank = 1")

	subQueryOrderProcessedDays := r.db.Model(&models.OrderDetail{}).
		Select(`
			order_details.order_id AS order_id,
			MAX(product_snapshots.process_day) AS max_processed_days
		`).
		Joins("JOIN product_snapshots ON product_snapshots.id = order_details.product_snapshot_id").
		Group("order_details.order_id")

	subQueryExpirationTime := r.db.Model(&models.Order{}).
		Select(`
			orders.id AS order_id,
			CASE
				WHEN orders.status = @pending_seller_status 
				THEN tbl_seller_confirmation.history_datetime + INTERVAL '1 minute' * (@interval_seller_confirmation_expiration)::INTEGER
				ELSE null
			END AS seller_confirmation_expiration_datetime,
			CASE
				WHEN orders.status = @on_process_status 
				THEN tbl_on_process.history_datetime + INTERVAL '1 day' * (tbl_processed_days.max_processed_days)::INTEGER
				ELSE null
			END AS on_process_expiration_datetime,
			CASE
				WHEN orders.status = @delivered_status 
				THEN tbl_delivered.history_datetime + INTERVAL '1 minute' * (@interval_auto_complete_delivered_expiration)::INTEGER
				ELSE null
			END AS delivered_expiration_datetime,
			CASE
				WHEN orders.status = @complained_status 
				THEN tbl_ongoing_last_refund_status.refund_request_datetime + INTERVAL '1 minute' * (@interval_auto_approve_refund_seller_confirmation_expiration)::INTEGER
				ELSE null
			END AS refund_seller_confirmation_expiration_datetime,
			CASE
				WHEN orders.status = @complained_status 
				THEN tbl_rejected_last_refund_status.rejected_datetime + INTERVAL '1 minute' * (@interval_auto_complete_order_rejected_refund_expiration)::INTEGER
				ELSE null
			END AS refund_rejected_auto_complete_order_expiration_datetime
		`,
			sql.Named("pending_seller_status", constants.PendingSellerStatus),
			sql.Named("on_process_status", constants.OnProcessStatus),
			sql.Named("delivered_status", constants.DeliveredStatus),
			sql.Named("complained_status", constants.ComplainedStatus),
			sql.Named("interval_seller_confirmation_expiration", sellerConfirmationExpirationInterval),
			sql.Named("interval_auto_complete_delivered_expiration", autoCompleteDeliveredExpirationInterval),
			sql.Named("interval_auto_approve_refund_seller_confirmation_expiration", autoApproveRefundSellerConfirmationExpirationInterval),
			sql.Named("interval_auto_complete_order_rejected_refund_expiration", autoCompleteOrderRejectedRefundExpirationInterval),
		).
		Joins("LEFT JOIN order_trackings tbl_seller_confirmation ON tbl_seller_confirmation.order_id = orders.id AND tbl_seller_confirmation.title = @payment_verified_title",
			sql.Named("payment_verified_title", constants.PaymentVerifiedTitle)).
		Joins("LEFT JOIN order_trackings tbl_on_process ON tbl_on_process.order_id = orders.id AND tbl_on_process.title = @order_processed_title",
			sql.Named("order_processed_title", constants.OrderProcessedTitle)).
		Joins("LEFT JOIN order_trackings tbl_delivered ON tbl_delivered.order_id = orders.id AND tbl_delivered.title = @delivered_title",
			sql.Named("delivered_title", constants.CronOrderTrackingTitleDelivered)).
		Joins("JOIN (?) tbl_processed_days ON tbl_processed_days.order_id = orders.id", subQueryOrderProcessedDays).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_rejected_last_refund_status ON orders.id = tbl_rejected_last_refund_status.order_id AND tbl_rejected_last_refund_status.status = @refund_rejected_status",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_rejected_status", constants.RefundStatusRejected)).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_ongoing_last_refund_status ON orders.id = tbl_ongoing_last_refund_status.order_id AND tbl_ongoing_last_refund_status.status = @refund_ongoing_status",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_ongoing_status", constants.RefundStatusOngoing))

	subQueryIsAllReviewed := r.db.Model(&models.Order{}).
		Select(`
			orders.id AS order_id, 
			ARRAY_AGG(product_reviews.id) AS product_review_id_array 
		`).
		Joins("JOIN order_details ON order_details.order_id = orders.id").
		Joins("LEFT JOIN product_reviews ON product_reviews.order_detail_id = order_details.id").
		Group("orders.id")

	queryResult := r.db.Model(&models.Order{}).
		Select(`
			orders.transaction_id AS transaction_id,
			orders.invoice_number AS invoice_number,
			orders.id AS order_id,
			orders.merchant_id AS merchant_id,
			merchants.name AS shop_name,
			merchants.photo AS shop_photo,
			users.id AS user_id,
			users.full_name AS user_full_name,
			courier_services.id AS courier_service_id,
			courier_services.description AS courier_description,
			transactions.transaction_datetime AS transaction_datetime,
			orders.status AS order_status,
			(orders.total_price_after_discount::float / transactions.total_price::float * transactions.total_price_after_discount) + orders.shipping_fee AS total_price,
			(
				CASE
					WHEN orders.status = @completed_status THEN tbl_completed_last_refund_status.refund_id
					ELSE null
				END
			) AS refund_id,
			(
				CASE
					WHEN orders.status = @complained_status THEN tbl_approved_last_refund_status.refund_id
					ELSE null
				END
			) AS approved_refund_id,
			(
				CASE
					WHEN orders.status = @complained_status THEN tbl_ongoing_last_refund_status.refund_id
					ELSE null
				END
			) AS ongoing_refund_id,
			orders.claim_status AS claim_status,
			NOT (SELECT bool_or(product_review_id IS NULL) FROM unnest(tbl_review_array.product_review_id_array) AS data(product_review_id)) AS is_all_reviewed,
			tbl_expiration_time.seller_confirmation_expiration_datetime AS seller_confirmation_expiration_datetime,
			tbl_expiration_time.on_process_expiration_datetime AS on_process_expiration_datetime,
			tbl_expiration_time.delivered_expiration_datetime AS delivered_expiration_datetime,
			tbl_expiration_time.refund_seller_confirmation_expiration_datetime AS refund_seller_confirmation_expiration_datetime,
			tbl_expiration_time.refund_rejected_auto_complete_order_expiration_datetime AS refund_rejected_auto_complete_order_expiration_datetime
		`,
			sql.Named("completed_status", constants.CompletedStatus),
			sql.Named("complained_status", constants.ComplainedStatus),
		).
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN users ON users.id = transactions.user_id").
		Joins("JOIN courier_services ON orders.courier_service_id = courier_services.id").
		Joins("LEFT JOIN (?) tbl_expiration_time ON tbl_expiration_time.order_id = orders.id", subQueryExpirationTime).
		Joins("JOIN (?) tbl_review_array ON tbl_review_array.order_id = orders.id", subQueryIsAllReviewed).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_completed_last_refund_status ON orders.id = tbl_completed_last_refund_status.order_id AND tbl_completed_last_refund_status.status = @refund_completed_status",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_completed_status", constants.RefundStatusCompleted),
		).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_approved_last_refund_status ON orders.id = tbl_approved_last_refund_status.order_id AND (tbl_approved_last_refund_status.status = @refund_approved_status OR tbl_approved_last_refund_status.status = @refund_approved_by_system_status)",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_approved_status", constants.RefundStatusApproved),
			sql.Named("refund_approved_by_system_status", constants.RefundStatusApprovedBySystem),
		).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_ongoing_last_refund_status ON orders.id = tbl_ongoing_last_refund_status.order_id AND tbl_ongoing_last_refund_status.status = @refund_ongoing_status",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_ongoing_status", constants.RefundStatusOngoing),
		).
		Where("transactions.user_id = ?", id).
		Where("LOWER(orders.status) != LOWER(?)", constants.PendingPaymentStatus)

	if query.FilterIsReviewed != "" && !strings.EqualFold(query.FilterIsReviewed, "all") {
		queryResult.Where("(NOT (SELECT bool_or(product_review_id IS NULL) FROM unnest(tbl_review_array.product_review_id_array) AS data(product_review_id)))::text = ?", query.FilterIsReviewed)
	}

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

	queryResult.
		Order("transactions.transaction_datetime DESC").
		Scopes(utils.Paginate(pagination, queryResult)).
		Find(&orderList)

	if queryResult.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return orderList, pagination, nil
}

func (r *orderRepository) OrderListGetAllByMerchantID(merchantID uint, query *externalDtos.OrderListQuery, pagination *utils.Pagination) ([]*dtos.OrderList, *utils.Pagination, error) {
	var orderList []*dtos.OrderList

	sellerConfirmationExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_SELLER_CONFIRMATION_EXPIRATION_INTERVAL"))
	if err != nil {
		sellerConfirmationExpirationInterval = constants.CronSellerConfirmationExpirationInterval
	}

	autoCompleteDeliveredExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_AUTO_COMPLETE_DELIVERED_EXPIRATION_INTERVAL"))
	if err != nil {
		autoCompleteDeliveredExpirationInterval = constants.CronAutoCompleteDeliveredExpirationInterval
	}

	autoCompleteOrderRejectedRefundExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_AUTO_COMPLETE_ORDER_REJECTED_REFUND_EXPIRATION_INTERVAL"))
	if err != nil {
		autoCompleteOrderRejectedRefundExpirationInterval = constants.CronAutoCompleteOrderRejectedRefundExpirationInterval
	}

	autoApproveRefundSellerConfirmationExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_AUTO_APPROVE_REFUND_SELLER_CONFIRMATION_EXPIRATION_INTERVAL"))
	if err != nil {
		autoApproveRefundSellerConfirmationExpirationInterval = constants.CronAutoApproveRefundSellerConfirmationExpirationInterval
	}

	subQueryRefundRank := r.db.Model(&models.Refund{}).
		Select(`
			id AS refund_id,
			order_id AS order_id,
			status AS status,
			refund_request_datetime AS refund_request_datetime,
			rejected_datetime AS rejected_datetime,
			RANK() OVER (PARTITION BY order_id ORDER BY refund_request_datetime DESC) AS refund_rank
		`)

	subQueryLastRefundStatus := r.db.Table("(?) AS tbl_refund_rank", subQueryRefundRank).
		Select(`
			tbl_refund_rank.refund_id AS refund_id,
			tbl_refund_rank.order_id AS order_id,
			tbl_refund_rank.status AS status,
			tbl_refund_rank.refund_request_datetime AS refund_request_datetime,
			tbl_refund_rank.rejected_datetime AS rejected_datetime
		`).
		Where("tbl_refund_rank.refund_rank = 1")

	subQueryOrderProcessedDays := r.db.Model(&models.OrderDetail{}).
		Select(`
			order_details.order_id AS order_id,
			MAX(product_snapshots.process_day) AS max_processed_days
		`).
		Joins("JOIN product_snapshots ON product_snapshots.id = order_details.product_snapshot_id").
		Group("order_details.order_id")

	subQueryExpirationTime := r.db.Model(&models.Order{}).
		Select(`
			orders.id AS order_id,
			CASE
				WHEN orders.status = @pending_seller_status 
				THEN tbl_seller_confirmation.history_datetime + INTERVAL '1 minute' * (@interval_seller_confirmation_expiration)::INTEGER
				ELSE null
			END AS seller_confirmation_expiration_datetime,
			CASE
				WHEN orders.status = @on_process_status 
				THEN tbl_on_process.history_datetime + INTERVAL '1 day' * (tbl_processed_days.max_processed_days)::INTEGER
				ELSE null
			END AS on_process_expiration_datetime,
			CASE
				WHEN orders.status = @delivered_status 
				THEN tbl_delivered.history_datetime + INTERVAL '1 minute' * (@interval_auto_complete_delivered_expiration)::INTEGER
				ELSE null
			END AS delivered_expiration_datetime,
			CASE
				WHEN orders.status = @complained_status 
				THEN tbl_ongoing_last_refund_status.refund_request_datetime + INTERVAL '1 minute' * (@interval_auto_approve_refund_seller_confirmation_expiration)::INTEGER
				ELSE null
			END AS refund_seller_confirmation_expiration_datetime,
			CASE
				WHEN orders.status = @complained_status 
				THEN tbl_rejected_last_refund_status.rejected_datetime + INTERVAL '1 minute' * (@interval_auto_complete_order_rejected_refund_expiration)::INTEGER
				ELSE null
			END AS refund_rejected_auto_complete_order_expiration_datetime
		`,
			sql.Named("pending_seller_status", constants.PendingSellerStatus),
			sql.Named("on_process_status", constants.OnProcessStatus),
			sql.Named("delivered_status", constants.DeliveredStatus),
			sql.Named("complained_status", constants.ComplainedStatus),
			sql.Named("interval_seller_confirmation_expiration", sellerConfirmationExpirationInterval),
			sql.Named("interval_auto_complete_delivered_expiration", autoCompleteDeliveredExpirationInterval),
			sql.Named("interval_auto_approve_refund_seller_confirmation_expiration", autoApproveRefundSellerConfirmationExpirationInterval),
			sql.Named("interval_auto_complete_order_rejected_refund_expiration", autoCompleteOrderRejectedRefundExpirationInterval),
		).
		Joins("LEFT JOIN order_trackings tbl_seller_confirmation ON tbl_seller_confirmation.order_id = orders.id AND tbl_seller_confirmation.title = @payment_verified_title",
			sql.Named("payment_verified_title", constants.PaymentVerifiedTitle)).
		Joins("LEFT JOIN order_trackings tbl_on_process ON tbl_on_process.order_id = orders.id AND tbl_on_process.title = @order_processed_title",
			sql.Named("order_processed_title", constants.OrderProcessedTitle)).
		Joins("LEFT JOIN order_trackings tbl_delivered ON tbl_delivered.order_id = orders.id AND tbl_delivered.title = @delivered_title",
			sql.Named("delivered_title", constants.CronOrderTrackingTitleDelivered)).
		Joins("JOIN (?) tbl_processed_days ON tbl_processed_days.order_id = orders.id", subQueryOrderProcessedDays).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_rejected_last_refund_status ON orders.id = tbl_rejected_last_refund_status.order_id AND tbl_rejected_last_refund_status.status = @refund_rejected_status",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_rejected_status", constants.RefundStatusRejected)).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_ongoing_last_refund_status ON orders.id = tbl_ongoing_last_refund_status.order_id AND tbl_ongoing_last_refund_status.status = @refund_ongoing_status",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_ongoing_status", constants.RefundStatusOngoing))

	subQueryIsAllReviewed := r.db.Model(&models.Order{}).
		Select(`
			orders.id AS order_id, 
			ARRAY_AGG(product_reviews.id) AS product_review_id_array 
		`).
		Joins("JOIN order_details ON order_details.order_id = orders.id").
		Joins("LEFT JOIN product_reviews ON product_reviews.order_detail_id = order_details.id").
		Group("orders.id")

	queryResult := r.db.Model(&models.Order{}).
		Select(`
			orders.transaction_id AS transaction_id,
			orders.invoice_number AS invoice_number,
			orders.id AS order_id,
			orders.merchant_id AS merchant_id,
			merchants.name AS shop_name,
			merchants.photo AS shop_photo,
			users.id AS user_id,
			users.full_name AS user_full_name,
			courier_services.id AS courier_service_id,
			courier_services.description AS courier_description,
			transactions.transaction_datetime AS transaction_datetime,
			orders.status AS order_status,
			orders.total_price_after_discount::float AS total_price,
			(
				CASE
					WHEN orders.status = @completed_status THEN tbl_completed_last_refund_status.refund_id
					ELSE null
				END
			) AS refund_id,
			(
				CASE
					WHEN orders.status = @complained_status THEN tbl_approved_last_refund_status.refund_id
					ELSE null
				END
			) AS approved_refund_id,
			(
				CASE
					WHEN orders.status = @complained_status THEN tbl_ongoing_last_refund_status.refund_id
					ELSE null
				END
			) AS ongoing_refund_id,
			orders.claim_status AS claim_status,
			NOT (SELECT bool_or(product_review_id IS NULL) FROM unnest(tbl_review_array.product_review_id_array) AS data(product_review_id)) AS is_all_reviewed,
			tbl_expiration_time.seller_confirmation_expiration_datetime AS seller_confirmation_expiration_datetime,
			tbl_expiration_time.on_process_expiration_datetime AS on_process_expiration_datetime,
			tbl_expiration_time.delivered_expiration_datetime AS delivered_expiration_datetime,
			tbl_expiration_time.refund_seller_confirmation_expiration_datetime AS refund_seller_confirmation_expiration_datetime,
			tbl_expiration_time.refund_rejected_auto_complete_order_expiration_datetime AS refund_rejected_auto_complete_order_expiration_datetime
		`,
			sql.Named("completed_status", constants.CompletedStatus),
			sql.Named("complained_status", constants.ComplainedStatus),
		).
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN users ON users.id = transactions.user_id").
		Joins("JOIN courier_services ON orders.courier_service_id = courier_services.id").
		Joins("LEFT JOIN (?) tbl_expiration_time ON tbl_expiration_time.order_id = orders.id", subQueryExpirationTime).
		Joins("JOIN (?) tbl_review_array ON tbl_review_array.order_id = orders.id", subQueryIsAllReviewed).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_completed_last_refund_status ON orders.id = tbl_completed_last_refund_status.order_id AND tbl_completed_last_refund_status.status = @refund_completed_status",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_completed_status", constants.RefundStatusCompleted),
		).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_approved_last_refund_status ON orders.id = tbl_approved_last_refund_status.order_id AND (tbl_approved_last_refund_status.status = @refund_approved_status OR tbl_approved_last_refund_status.status = @refund_approved_by_system_status)",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_approved_status", constants.RefundStatusApproved),
			sql.Named("refund_approved_by_system_status", constants.RefundStatusApprovedBySystem),
		).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_ongoing_last_refund_status ON orders.id = tbl_ongoing_last_refund_status.order_id AND tbl_ongoing_last_refund_status.status = @refund_ongoing_status",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_ongoing_status", constants.RefundStatusOngoing),
		).
		Where("orders.merchant_id = ?", merchantID)
		// Where("LOWER(orders.status) != LOWER(?)", constants.PendingPaymentStatus)

	if query.FilterIsReviewed != "" && !strings.EqualFold(query.FilterIsReviewed, "all") {
		queryResult.Where("(NOT (SELECT bool_or(product_review_id IS NULL) FROM unnest(tbl_review_array.product_review_id_array) AS data(product_review_id)))::text = ?", query.FilterIsReviewed)
	}

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

	queryResult.
		Order("transactions.transaction_datetime DESC").
		Scopes(utils.Paginate(pagination, queryResult)).
		Find(&orderList)

	if queryResult.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return orderList, pagination, nil
}

func (r *orderRepository) OrderListDetailsGetByID(id uuid.UUID) (*dtos.OrderListDetails, error) {
	var orderListDetails *dtos.OrderListDetails

	sellerConfirmationExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_SELLER_CONFIRMATION_EXPIRATION_INTERVAL"))
	if err != nil {
		sellerConfirmationExpirationInterval = constants.CronSellerConfirmationExpirationInterval
	}

	autoCompleteDeliveredExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_AUTO_COMPLETE_DELIVERED_EXPIRATION_INTERVAL"))
	if err != nil {
		autoCompleteDeliveredExpirationInterval = constants.CronAutoCompleteDeliveredExpirationInterval
	}

	autoCompleteOrderRejectedRefundExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_AUTO_COMPLETE_ORDER_REJECTED_REFUND_EXPIRATION_INTERVAL"))
	if err != nil {
		autoCompleteOrderRejectedRefundExpirationInterval = constants.CronAutoCompleteOrderRejectedRefundExpirationInterval
	}

	autoApproveRefundSellerConfirmationExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_AUTO_APPROVE_REFUND_SELLER_CONFIRMATION_EXPIRATION_INTERVAL"))
	if err != nil {
		autoApproveRefundSellerConfirmationExpirationInterval = constants.CronAutoApproveRefundSellerConfirmationExpirationInterval
	}

	subQueryRefundRank := r.db.Model(&models.Refund{}).
		Select(`
			id AS refund_id,
			order_id AS order_id,
			status AS status,
			refund_request_datetime AS refund_request_datetime,
			rejected_datetime AS rejected_datetime,
			RANK() OVER (PARTITION BY order_id ORDER BY refund_request_datetime DESC) AS refund_rank
		`)

	subQueryLastRefundStatus := r.db.Table("(?) AS tbl_refund_rank", subQueryRefundRank).
		Select(`
			tbl_refund_rank.refund_id AS refund_id,
			tbl_refund_rank.order_id AS order_id,
			tbl_refund_rank.status AS status,
			tbl_refund_rank.refund_request_datetime AS refund_request_datetime,
			tbl_refund_rank.rejected_datetime AS rejected_datetime
		`).
		Where("tbl_refund_rank.refund_rank = 1")

	subQueryOrderProcessedDays := r.db.Model(&models.OrderDetail{}).
		Select(`
			order_details.order_id AS order_id,
			MAX(product_snapshots.process_day) AS max_processed_days
		`).
		Joins("JOIN product_snapshots ON product_snapshots.id = order_details.product_snapshot_id").
		Group("order_details.order_id")

	subQueryExpirationTime := r.db.Model(&models.Order{}).
		Select(`
			orders.id AS order_id,
			CASE
				WHEN orders.status = @pending_seller_status 
				THEN tbl_seller_confirmation.history_datetime + INTERVAL '1 minute' * (@interval_seller_confirmation_expiration)::INTEGER
				ELSE null
			END AS seller_confirmation_expiration_datetime,
			CASE
				WHEN orders.status = @on_process_status 
				THEN tbl_on_process.history_datetime + INTERVAL '1 day' * (tbl_processed_days.max_processed_days)::INTEGER
				ELSE null
			END AS on_process_expiration_datetime,
			CASE
				WHEN orders.status = @delivered_status 
				THEN tbl_delivered.history_datetime + INTERVAL '1 minute' * (@interval_auto_complete_delivered_expiration)::INTEGER
				ELSE null
			END AS delivered_expiration_datetime,
			CASE
				WHEN orders.status = @complained_status 
				THEN tbl_ongoing_last_refund_status.refund_request_datetime + INTERVAL '1 minute' * (@interval_auto_approve_refund_seller_confirmation_expiration)::INTEGER
				ELSE null
			END AS refund_seller_confirmation_expiration_datetime,
			CASE
				WHEN orders.status = @complained_status 
				THEN tbl_rejected_last_refund_status.rejected_datetime + INTERVAL '1 minute' * (@interval_auto_complete_order_rejected_refund_expiration)::INTEGER
				ELSE null
			END AS refund_rejected_auto_complete_order_expiration_datetime
		`,
			sql.Named("pending_seller_status", constants.PendingSellerStatus),
			sql.Named("on_process_status", constants.OnProcessStatus),
			sql.Named("delivered_status", constants.DeliveredStatus),
			sql.Named("complained_status", constants.ComplainedStatus),
			sql.Named("interval_seller_confirmation_expiration", sellerConfirmationExpirationInterval),
			sql.Named("interval_auto_complete_delivered_expiration", autoCompleteDeliveredExpirationInterval),
			sql.Named("interval_auto_approve_refund_seller_confirmation_expiration", autoApproveRefundSellerConfirmationExpirationInterval),
			sql.Named("interval_auto_complete_order_rejected_refund_expiration", autoCompleteOrderRejectedRefundExpirationInterval),
		).
		Joins("LEFT JOIN order_trackings tbl_seller_confirmation ON tbl_seller_confirmation.order_id = orders.id AND tbl_seller_confirmation.title = @payment_verified_title",
			sql.Named("payment_verified_title", constants.PaymentVerifiedTitle)).
		Joins("LEFT JOIN order_trackings tbl_on_process ON tbl_on_process.order_id = orders.id AND tbl_on_process.title = @order_processed_title",
			sql.Named("order_processed_title", constants.OrderProcessedTitle)).
		Joins("LEFT JOIN order_trackings tbl_delivered ON tbl_delivered.order_id = orders.id AND tbl_delivered.title = @delivered_title",
			sql.Named("delivered_title", constants.CronOrderTrackingTitleDelivered)).
		Joins("JOIN (?) tbl_processed_days ON tbl_processed_days.order_id = orders.id", subQueryOrderProcessedDays).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_rejected_last_refund_status ON orders.id = tbl_rejected_last_refund_status.order_id AND tbl_rejected_last_refund_status.status = @refund_rejected_status",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_rejected_status", constants.RefundStatusRejected)).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_ongoing_last_refund_status ON orders.id = tbl_ongoing_last_refund_status.order_id AND tbl_ongoing_last_refund_status.status = @refund_ongoing_status",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_ongoing_status", constants.RefundStatusOngoing))

	subQueryIsAllReviewed := r.db.Model(&models.Order{}).
		Select(`
		orders.id AS order_id, 
		ARRAY_AGG(product_reviews.id) AS product_review_id_array 
	`).
		Joins("JOIN order_details ON order_details.order_id = orders.id").
		Joins("LEFT JOIN product_reviews ON product_reviews.order_detail_id = order_details.id").
		Group("orders.id")

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
			merchants.photo AS shop_photo,
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
			(
				CASE
					WHEN orders.status = @completed_status THEN tbl_completed_last_refund_status.refund_id
					ELSE null
				END
			) AS refund_id,
			(
				CASE
					WHEN orders.status = @complained_status THEN tbl_approved_last_refund_status.refund_id
					ELSE null
				END
			) AS approved_refund_id,
			(
				CASE
					WHEN orders.status = @complained_status THEN tbl_ongoing_last_refund_status.refund_id
					ELSE null
				END
			) AS ongoing_refund_id,
			orders.claim_status AS claim_status,
			NOT (SELECT bool_or(product_review_id IS NULL) FROM unnest(tbl_review_array.product_review_id_array) AS data(product_review_id)) AS is_all_reviewed,
			tbl_expiration_time.seller_confirmation_expiration_datetime AS seller_confirmation_expiration_datetime,
			tbl_expiration_time.on_process_expiration_datetime AS on_process_expiration_datetime,
			tbl_expiration_time.delivered_expiration_datetime AS delivered_expiration_datetime,
			tbl_expiration_time.refund_seller_confirmation_expiration_datetime AS refund_seller_confirmation_expiration_datetime,
			tbl_expiration_time.refund_rejected_auto_complete_order_expiration_datetime AS refund_rejected_auto_complete_order_expiration_datetime
		`,
			subQueryTotalPriceBeforeDiscount,
			subQueryMerchantPromotion,
			subQueryMerchantVoucher,
			subQueryMarketplaceVoucher,
			sql.Named("completed_status", constants.CompletedStatus),
			sql.Named("complained_status", constants.ComplainedStatus),
		).
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN users ON users.id = transactions.user_id").
		Joins("JOIN payment_methods ON transactions.payment_method_id = payment_methods.id").
		Joins("JOIN courier_services ON orders.courier_service_id = courier_services.id").
		Joins("LEFT JOIN (?) tbl_expiration_time ON tbl_expiration_time.order_id = orders.id", subQueryExpirationTime).
		Joins("JOIN (?) tbl_review_array ON tbl_review_array.order_id = orders.id", subQueryIsAllReviewed).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_completed_last_refund_status ON orders.id = tbl_completed_last_refund_status.order_id AND tbl_completed_last_refund_status.status = @refund_completed_status",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_completed_status", constants.RefundStatusCompleted),
		).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_approved_last_refund_status ON orders.id = tbl_approved_last_refund_status.order_id AND (tbl_approved_last_refund_status.status = @refund_approved_status OR tbl_approved_last_refund_status.status = @refund_approved_by_system_status)",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_approved_status", constants.RefundStatusApproved),
			sql.Named("refund_approved_by_system_status", constants.RefundStatusApprovedBySystem),
		).
		Joins("LEFT JOIN (@tbl_last_refund_status) tbl_ongoing_last_refund_status ON orders.id = tbl_ongoing_last_refund_status.order_id AND tbl_ongoing_last_refund_status.status = @refund_ongoing_status",
			sql.Named("tbl_last_refund_status", subQueryLastRefundStatus),
			sql.Named("refund_ongoing_status", constants.RefundStatusOngoing),
		).
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
	queryResult := r.db.Model(&models.Order{}).
		Preload("Transaction").
		First(&model, id)
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

func (r *orderRepository) GetByTransactionID(transactionID uuid.UUID) (models []*models.Order, err error) {
	queryResult := r.db.
		Where("transaction_id = ?", transactionID).
		Find(&models)
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

func (r *orderRepository) PendingPaymentGetAll() (pendingOrders []*models.Order, err error) {
	queryResult := r.db.Model(&models.Order{}).
		Preload("Transaction").
		Where("orders.status = ?", constants.PendingPaymentStatus).
		Find(&pendingOrders)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *orderRepository) GetMaxProcessedDaysByOrderID(orderID uuid.UUID) (maxProcessedDays uint, err error) {
	queryResult := r.db.Model(&models.OrderDetail{}).
		Select(`
			MAX(product_snapshots.process_day) AS max_processed_days
		`).
		Where("order_details.order_id = ?", orderID).
		Joins("JOIN product_snapshots ON product_snapshots.id = order_details.product_snapshot_id").
		Group("order_details.order_id").
		Scan(&maxProcessedDays)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}
