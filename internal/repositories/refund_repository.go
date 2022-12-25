package repositories

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/constants/refund_media_constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefundRepository interface {
	Create(newModel *models.Refund) (model *models.Refund, err error)
	GetAllDetailsByBuyer(userID uuid.UUID, query *externalDtos.GetAllRefundQuery, pagination *utils.Pagination) ([]*dtos.RefundDetails, *utils.Pagination, error)
	GetDetailsByID(refundID uuid.UUID) (*dtos.RefundDetails, error)
	GetAllDetailsByMerchant(merchantID uint, query *externalDtos.GetAllRefundQuery, pagination *utils.Pagination) ([]*dtos.RefundDetails, *utils.Pagination, error)
	GetAllDetails(query *externalDtos.GetAllRefundQuery, pagination *utils.Pagination) ([]*dtos.RefundDetails, *utils.Pagination, error)
	Update(id uuid.UUID, newModel *models.Refund) (model *models.Refund, err error)
	GetByID(id uuid.UUID) (model *models.Refund, err error)
	GetDetailsByOrderIDAndStatus(orderID uuid.UUID, status string) (refundDetails []*dtos.RefundDetails, err error)
	GetByOrderID(orderID uuid.UUID) (model []*models.Refund, err error)
	GetComplainedOrderLastRefund() (model []*models.Refund, err error)
}

type refundRepository struct {
	db *gorm.DB
}

type RefundRepositoryConfig struct {
	DB *gorm.DB
}

func NewRefundRepository(c *RefundRepositoryConfig) RefundRepository {
	return &refundRepository{
		db: c.DB,
	}
}

func (r *refundRepository) Create(newModel *models.Refund) (model *models.Refund, err error) {
	result := r.db.
		Where("order_id = ?", newModel.OrderID).
		Where("status IN (?, ?, ?, ?)", constants.RefundStatusOngoing, constants.RefundStatusApproved, constants.RefundStatusApprovedBySystem, constants.RefundStatusCompleted).
		FirstOrCreate(&newModel)

	if result.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	if result.RowsAffected == 0 {
		err = e.ErrRefundAlreadyExist
		return
	}

	model = newModel

	return
}

func (r *refundRepository) GetAllDetailsByBuyer(userID uuid.UUID, query *externalDtos.GetAllRefundQuery, pagination *utils.Pagination) ([]*dtos.RefundDetails, *utils.Pagination, error) {
	var refundDetails []*dtos.RefundDetails

	autoApproveRefundSellerConfirmationExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_AUTO_APPROVE_REFUND_SELLER_CONFIRMATION_EXPIRATION_INTERVAL"))
	if err != nil {
		autoApproveRefundSellerConfirmationExpirationInterval = 2880
	}

	subQueryExpirationTime := r.db.Model(&models.Refund{}).
		Select(`
			refunds.id AS refund_id,
			CASE
				WHEN refunds.status = @refund_ongoing_status
				THEN refunds.refund_request_datetime + INTERVAL '1 minute' * (@interval_auto_approve_refund_seller_confirmation_expiration)::INTEGER
				ELSE null
			END AS refund_seller_confirmation_expiration_datetime
		`,
			sql.Named("refund_ongoing_status", constants.RefundStatusOngoing),
			sql.Named("interval_auto_approve_refund_seller_confirmation_expiration", autoApproveRefundSellerConfirmationExpirationInterval),
		)

	subQueryRefundMedia := r.db.Model(&models.RefundMedia{}).
		Select(`
			refund_id AS refund_id,
			type AS type,
			ARRAY_AGG(media_url) AS media_url
		`).
		Group("refund_id, type")

	subQueryRefundDetails := r.db.Model(&models.RefundDetail{}).
		Select(`
				refund_id AS refund_id,
				SUM(amount_user_get) AS total_user_refund_amount,
				SUM(amount_merchant_lost) AS total_merchant_lost_amount
			`).
		Group("refund_id")

	result := r.db.Model(&models.Refund{}).
		Select(`
			refunds.id AS refund_id,
			refunds.order_id AS order_id,
			orders.invoice_number AS invoice_number,
			orders.merchant_id AS merchant_id,
			merchants.name AS shop_name,
			transactions.user_id AS user_id,
			users.full_name AS user_full_name,
			refunds.refund_request_datetime AS refund_request_datetime,
			refunds.rejected_datetime AS rejected_datetime,
			refunds.status AS refund_status,
			refunds.reason AS refund_reason,
			tbl_refund_images.media_url AS images,
			tbl_refund_videos.media_url AS videos,
			tbl_refund_details.total_user_refund_amount AS total_user_refund_amount,
			tbl_refund_details.total_merchant_lost_amount AS total_merchant_lost_amount,
			tbl_expiration_time.refund_seller_confirmation_expiration_datetime AS refund_seller_confirmation_expiration_datetime
		`).
		Joins("JOIN orders ON refunds.order_id = orders.id").
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN users ON transactions.user_id = users.id").
		Joins("LEFT JOIN (?) tbl_refund_images ON refunds.id = tbl_refund_images.refund_id AND tbl_refund_images.type = ?", subQueryRefundMedia, refund_media_constants.TypeImage).
		Joins("LEFT JOIN (?) tbl_refund_videos ON refunds.id = tbl_refund_videos.refund_id AND tbl_refund_videos.type = ?", subQueryRefundMedia, refund_media_constants.TypeVideo).
		Joins("LEFT JOIN (?) tbl_refund_details ON refunds.id = tbl_refund_details.refund_id", subQueryRefundDetails).
		Joins("LEFT JOIN (?) tbl_expiration_time ON refunds.id = tbl_expiration_time.refund_id", subQueryExpirationTime).
		Where("transactions.user_id = ?", userID)

	if query.SearchInvoice != "" {
		result.Where("orders.invoice_number = ?", query.SearchInvoice)
	}

	if query.FilterGroupedStatus != "" && query.FilterGroupedStatus != constants.RefundFilterAll {
		if query.FilterGroupedStatus == constants.RefundGroupedStatusProcessed {
			result.Where("refunds.status IN (?, ?, ?, ?)", constants.RefundStatusApproved, constants.RefundStatusApprovedBySystem, constants.RefundStatusCompleted, constants.RefundStatusRejected)
		}
	}

	if query.FilterStatus != "" && query.FilterStatus != constants.RefundFilterAll {
		result.Where("refunds.status = ?", query.FilterStatus)
	}

	if query.FilterStartDate != nil && query.FilterEndDate != nil {
		result.Where("refunds.refund_request_datetime::date BETWEEN ?::date AND ?::date", query.FilterStartDate, query.FilterEndDate)
	}

	result.
		Order(fmt.Sprintf("%s %s", query.SortBy, query.Sort)).
		Scopes(utils.Paginate(pagination, result)).
		Find(&refundDetails)

	if result.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return refundDetails, pagination, nil
}

func (r *refundRepository) GetDetailsByID(refundID uuid.UUID) (*dtos.RefundDetails, error) {
	var refundDetails *dtos.RefundDetails

	autoApproveRefundSellerConfirmationExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_AUTO_APPROVE_REFUND_SELLER_CONFIRMATION_EXPIRATION_INTERVAL"))
	if err != nil {
		autoApproveRefundSellerConfirmationExpirationInterval = 2880
	}

	subQueryExpirationTime := r.db.Model(&models.Refund{}).
		Select(`
			refunds.id AS refund_id,
			CASE
				WHEN refunds.status = @refund_ongoing_status
				THEN refunds.refund_request_datetime + INTERVAL '1 minute' * (@interval_auto_approve_refund_seller_confirmation_expiration)::INTEGER
				ELSE null
			END AS refund_seller_confirmation_expiration_datetime
		`,
			sql.Named("refund_ongoing_status", constants.RefundStatusOngoing),
			sql.Named("interval_auto_approve_refund_seller_confirmation_expiration", autoApproveRefundSellerConfirmationExpirationInterval),
		)

	subQueryRefundMedia := r.db.Model(&models.RefundMedia{}).
		Select(`
			refund_id AS refund_id,
			type AS type,
			ARRAY_AGG(media_url) AS media_url
		`).
		Group("refund_id, type")

	subQueryRefundDetails := r.db.Model(&models.RefundDetail{}).
		Select(`
				refund_id AS refund_id,
				SUM(amount_user_get) AS total_user_refund_amount,
				SUM(amount_merchant_lost) AS total_merchant_lost_amount
			`).
		Group("refund_id")

	result := r.db.Model(&models.Refund{}).
		Select(`
			refunds.id AS refund_id,
			refunds.order_id AS order_id,
			orders.invoice_number AS invoice_number,
			orders.merchant_id AS merchant_id,
			merchants.name AS shop_name,
			transactions.user_id AS user_id,
			users.full_name AS user_full_name,
			refunds.refund_request_datetime AS refund_request_datetime,
			refunds.rejected_datetime AS rejected_datetime,
			refunds.status AS refund_status,
			refunds.reason AS refund_reason,
			tbl_refund_images.media_url AS images,
			tbl_refund_videos.media_url AS videos,
			tbl_refund_details.total_user_refund_amount AS total_user_refund_amount,
			tbl_refund_details.total_merchant_lost_amount AS total_merchant_lost_amount,
			tbl_expiration_time.refund_seller_confirmation_expiration_datetime AS refund_seller_confirmation_expiration_datetime
		`).
		Joins("JOIN orders ON refunds.order_id = orders.id").
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN users ON transactions.user_id = users.id").
		Joins("LEFT JOIN (?) tbl_refund_images ON refunds.id = tbl_refund_images.refund_id AND tbl_refund_images.type = ?", subQueryRefundMedia, refund_media_constants.TypeImage).
		Joins("LEFT JOIN (?) tbl_refund_videos ON refunds.id = tbl_refund_videos.refund_id AND tbl_refund_videos.type = ?", subQueryRefundMedia, refund_media_constants.TypeVideo).
		Joins("LEFT JOIN (?) tbl_refund_details ON refunds.id = tbl_refund_details.refund_id", subQueryRefundDetails).
		Joins("LEFT JOIN (?) tbl_expiration_time ON refunds.id = tbl_expiration_time.refund_id", subQueryExpirationTime).
		Where("refunds.id = ?", refundID).
		First(&refundDetails)

	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, e.ErrDataNotFound
		}

		return nil, e.ErrDatabaseOperationFailed
	}

	return refundDetails, nil
}

func (r *refundRepository) GetAllDetailsByMerchant(merchantID uint, query *externalDtos.GetAllRefundQuery, pagination *utils.Pagination) ([]*dtos.RefundDetails, *utils.Pagination, error) {
	var refundDetails []*dtos.RefundDetails

	autoApproveRefundSellerConfirmationExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_AUTO_APPROVE_REFUND_SELLER_CONFIRMATION_EXPIRATION_INTERVAL"))
	if err != nil {
		autoApproveRefundSellerConfirmationExpirationInterval = 2880
	}

	subQueryExpirationTime := r.db.Model(&models.Refund{}).
		Select(`
			refunds.id AS refund_id,
			CASE
				WHEN refunds.status = @refund_ongoing_status
				THEN refunds.refund_request_datetime + INTERVAL '1 minute' * (@interval_auto_approve_refund_seller_confirmation_expiration)::INTEGER
				ELSE null
			END AS refund_seller_confirmation_expiration_datetime
		`,
			sql.Named("refund_ongoing_status", constants.RefundStatusOngoing),
			sql.Named("interval_auto_approve_refund_seller_confirmation_expiration", autoApproveRefundSellerConfirmationExpirationInterval),
		)

	subQueryRefundMedia := r.db.Model(&models.RefundMedia{}).
		Select(`
			refund_id AS refund_id,
			type AS type,
			ARRAY_AGG(media_url) AS media_url
		`).
		Group("refund_id, type")

	subQueryRefundDetails := r.db.Model(&models.RefundDetail{}).
		Select(`
				refund_id AS refund_id,
				SUM(amount_user_get) AS total_user_refund_amount,
				SUM(amount_merchant_lost) AS total_merchant_lost_amount
			`).
		Group("refund_id")

	result := r.db.Model(&models.Refund{}).
		Select(`
			refunds.id AS refund_id,
			refunds.order_id AS order_id,
			orders.invoice_number AS invoice_number,
			orders.merchant_id AS merchant_id,
			merchants.name AS shop_name,
			transactions.user_id AS user_id,
			users.full_name AS user_full_name,
			refunds.refund_request_datetime AS refund_request_datetime,
			refunds.rejected_datetime AS rejected_datetime,
			refunds.status AS refund_status,
			refunds.reason AS refund_reason,
			tbl_refund_images.media_url AS images,
			tbl_refund_videos.media_url AS videos,
			tbl_refund_details.total_user_refund_amount AS total_user_refund_amount,
			tbl_refund_details.total_merchant_lost_amount AS total_merchant_lost_amount,
			tbl_expiration_time.refund_seller_confirmation_expiration_datetime AS refund_seller_confirmation_expiration_datetime
		`).
		Joins("JOIN orders ON refunds.order_id = orders.id").
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN users ON transactions.user_id = users.id").
		Joins("LEFT JOIN (?) tbl_refund_images ON refunds.id = tbl_refund_images.refund_id AND tbl_refund_images.type = ?", subQueryRefundMedia, refund_media_constants.TypeImage).
		Joins("LEFT JOIN (?) tbl_refund_videos ON refunds.id = tbl_refund_videos.refund_id AND tbl_refund_videos.type = ?", subQueryRefundMedia, refund_media_constants.TypeVideo).
		Joins("LEFT JOIN (?) tbl_refund_details ON refunds.id = tbl_refund_details.refund_id", subQueryRefundDetails).
		Joins("LEFT JOIN (?) tbl_expiration_time ON refunds.id = tbl_expiration_time.refund_id", subQueryExpirationTime).
		Where("orders.merchant_id = ?", merchantID)

	if query.SearchInvoice != "" {
		result.Where("orders.invoice_number = ?", query.SearchInvoice)
	}

	if query.FilterGroupedStatus != "" && query.FilterGroupedStatus != constants.RefundFilterAll {
		if query.FilterGroupedStatus == constants.RefundGroupedStatusProcessed {
			result.Where("refunds.status IN (?, ?, ?, ?)", constants.RefundStatusApproved, constants.RefundStatusApprovedBySystem, constants.RefundStatusCompleted, constants.RefundStatusRejected)
		}
	}

	if query.FilterStatus != "" && query.FilterStatus != constants.RefundFilterAll {
		result.Where("refunds.status = ?", query.FilterStatus)
	}

	if query.FilterStartDate != nil && query.FilterEndDate != nil {
		result.Where("refunds.refund_request_datetime::date BETWEEN ?::date AND ?::date", query.FilterStartDate, query.FilterEndDate)
	}

	result.
		Order(fmt.Sprintf("%s %s", query.SortBy, query.Sort)).
		Scopes(utils.Paginate(pagination, result)).
		Find(&refundDetails)

	if result.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return refundDetails, pagination, nil
}

func (r *refundRepository) GetAllDetails(query *externalDtos.GetAllRefundQuery, pagination *utils.Pagination) ([]*dtos.RefundDetails, *utils.Pagination, error) {
	var refundDetails []*dtos.RefundDetails

	autoApproveRefundSellerConfirmationExpirationInterval, err := strconv.Atoi(os.Getenv("CRON_AUTO_APPROVE_REFUND_SELLER_CONFIRMATION_EXPIRATION_INTERVAL"))
	if err != nil {
		autoApproveRefundSellerConfirmationExpirationInterval = 2880
	}

	subQueryExpirationTime := r.db.Model(&models.Refund{}).
		Select(`
			refunds.id AS refund_id,
			CASE
				WHEN refunds.status = @refund_ongoing_status
				THEN refunds.refund_request_datetime + INTERVAL '1 minute' * (@interval_auto_approve_refund_seller_confirmation_expiration)::INTEGER
				ELSE null
			END AS refund_seller_confirmation_expiration_datetime
		`,
			sql.Named("refund_ongoing_status", constants.RefundStatusOngoing),
			sql.Named("interval_auto_approve_refund_seller_confirmation_expiration", autoApproveRefundSellerConfirmationExpirationInterval),
		)

	subQueryRefundMedia := r.db.Model(&models.RefundMedia{}).
		Select(`
			refund_id AS refund_id,
			type AS type,
			ARRAY_AGG(media_url) AS media_url
		`).
		Group("refund_id, type")

	subQueryRefundDetails := r.db.Model(&models.RefundDetail{}).
		Select(`
				refund_id AS refund_id,
				SUM(amount_user_get) AS total_user_refund_amount,
				SUM(amount_merchant_lost) AS total_merchant_lost_amount
			`).
		Group("refund_id")

	result := r.db.Model(&models.Refund{}).
		Select(`
			refunds.id AS refund_id,
			refunds.order_id AS order_id,
			orders.invoice_number AS invoice_number,
			orders.merchant_id AS merchant_id,
			merchants.name AS shop_name,
			transactions.user_id AS user_id,
			users.full_name AS user_full_name,
			refunds.refund_request_datetime AS refund_request_datetime,
			refunds.rejected_datetime AS rejected_datetime,
			refunds.status AS refund_status,
			refunds.reason AS refund_reason,
			tbl_refund_images.media_url AS images,
			tbl_refund_videos.media_url AS videos,
			tbl_refund_details.total_user_refund_amount AS total_user_refund_amount,
			tbl_refund_details.total_merchant_lost_amount AS total_merchant_lost_amount,
			tbl_expiration_time.refund_seller_confirmation_expiration_datetime AS refund_seller_confirmation_expiration_datetime
		`).
		Joins("JOIN orders ON refunds.order_id = orders.id").
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN users ON transactions.user_id = users.id").
		Joins("LEFT JOIN (?) tbl_refund_images ON refunds.id = tbl_refund_images.refund_id AND tbl_refund_images.type = ?", subQueryRefundMedia, refund_media_constants.TypeImage).
		Joins("LEFT JOIN (?) tbl_refund_videos ON refunds.id = tbl_refund_videos.refund_id AND tbl_refund_videos.type = ?", subQueryRefundMedia, refund_media_constants.TypeVideo).
		Joins("LEFT JOIN (?) tbl_refund_details ON refunds.id = tbl_refund_details.refund_id", subQueryRefundDetails).
		Joins("LEFT JOIN (?) tbl_expiration_time ON refunds.id = tbl_expiration_time.refund_id", subQueryExpirationTime)

	if query.SearchInvoice != "" {
		result.Where("orders.invoice_number = ?", query.SearchInvoice)
	}

	if query.FilterGroupedStatus != "" && query.FilterGroupedStatus != constants.RefundFilterAll {
		if query.FilterGroupedStatus == constants.RefundGroupedStatusProcessed {
			result.Where("refunds.status IN (?, ?, ?, ?)", constants.RefundStatusApproved, constants.RefundStatusApprovedBySystem, constants.RefundStatusCompleted, constants.RefundStatusRejected)
		}
	}

	if query.FilterStatus != "" && query.FilterStatus != constants.RefundFilterAll {
		result.Where("refunds.status = ?", query.FilterStatus)
	}

	if query.FilterStartDate != nil && query.FilterEndDate != nil {
		result.Where("refunds.refund_request_datetime::date BETWEEN ?::date AND ?::date", query.FilterStartDate, query.FilterEndDate)
	}

	result.
		Order(fmt.Sprintf("%s %s", query.SortBy, query.Sort)).
		Scopes(utils.Paginate(pagination, result)).
		Find(&refundDetails)

	if result.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return refundDetails, pagination, nil
}

func (r *refundRepository) Update(id uuid.UUID, newModel *models.Refund) (model *models.Refund, err error) {
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

func (r *refundRepository) GetByID(id uuid.UUID) (model *models.Refund, err error) {
	queryResult := r.db.Model(&models.Refund{}).
		Preload("Order").
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

func (r *refundRepository) GetDetailsByOrderIDAndStatus(orderID uuid.UUID, status string) (refundDetails []*dtos.RefundDetails, err error) {
	subQueryRefundMedia := r.db.Model(&models.RefundMedia{}).
		Select(`
				refund_id AS refund_id,
				type AS type,
				ARRAY_AGG(media_url) AS media_url
			`).
		Group("refund_id, type")

	subQueryRefundDetails := r.db.Model(&models.RefundDetail{}).
		Select(`
					refund_id AS refund_id,
					SUM(amount_user_get) AS total_user_refund_amount,
					SUM(amount_merchant_lost) AS total_merchant_lost_amount
				`).
		Group("refund_id")

	result := r.db.Model(&models.Refund{}).
		Select(`
				refunds.id AS refund_id,
				refunds.order_id AS order_id,
				orders.invoice_number AS invoice_number,
				orders.merchant_id AS merchant_id,
				merchants.name AS shop_name,
				transactions.user_id AS user_id,
				users.full_name AS user_full_name,
				refunds.refund_request_datetime AS refund_request_datetime,
				refunds.rejected_datetime AS rejected_datetime,
				refunds.status AS refund_status,
				refunds.reason AS refund_reason,
				tbl_refund_images.media_url AS images,
				tbl_refund_videos.media_url AS videos,
				tbl_refund_details.total_user_refund_amount AS total_user_refund_amount,
				tbl_refund_details.total_merchant_lost_amount AS total_merchant_lost_amount
			`).
		Joins("JOIN orders ON refunds.order_id = orders.id").
		Joins("JOIN merchants ON orders.merchant_id = merchants.id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("JOIN users ON transactions.user_id = users.id").
		Joins("LEFT JOIN (?) tbl_refund_images ON refunds.id = tbl_refund_images.refund_id AND tbl_refund_images.type = ?", subQueryRefundMedia, refund_media_constants.TypeImage).
		Joins("LEFT JOIN (?) tbl_refund_videos ON refunds.id = tbl_refund_videos.refund_id AND tbl_refund_videos.type = ?", subQueryRefundMedia, refund_media_constants.TypeVideo).
		Joins("LEFT JOIN (?) tbl_refund_details ON refunds.id = tbl_refund_details.refund_id", subQueryRefundDetails).
		Where("refunds.status = ?", status).
		Where("refunds.order_id = ?", orderID).
		Find(&refundDetails)

	if result.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *refundRepository) GetByOrderID(orderID uuid.UUID) (model []*models.Refund, err error) {

	result := r.db.Model(&models.Refund{}).
		Where("order_id = ?", orderID).
		Find(&model)

	if result.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *refundRepository) GetComplainedOrderLastRefund() (model []*models.Refund, err error) {
	subQueryRefundRank := r.db.Model(&models.Refund{}).
		Select(`
			id AS refund_id,
			status AS status,
			refund_request_datetime AS refund_request_datetime,
			rejected_datetime AS rejected_datetime,
			RANK() OVER (PARTITION BY order_id ORDER BY refund_request_datetime DESC) AS refund_rank
		`)

	result := r.db.Model(&models.Refund{}).
		Joins("JOIN (?) tbl_last_refund ON refunds.id = tbl_last_refund.refund_id AND tbl_last_refund.refund_rank = 1", subQueryRefundRank).
		Joins("JOIN orders ON refunds.order_id = orders.id AND orders.status = ?", constants.ComplainedStatus).
		Find(&model)

	if result.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}
