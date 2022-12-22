package dtos

import (
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type OrderDetailRefund struct {
	OrderDetailID uuid.UUID `json:"order_detail_id" binding:"required"`
	Quantity      uint      `json:"quantity" binding:"required"`
}

type CreateRefundRequest struct {
	OrderDetail []*OrderDetailRefund  `json:"order_detail" binding:"required"`
	Reason      string                `json:"reason" binding:"required"`
	Photo1      *multipart.FileHeader `form:"photo_1"`
	Photo2      *multipart.FileHeader `form:"photo_2"`
	Photo3      *multipart.FileHeader `form:"photo_3"`
	Photo4      *multipart.FileHeader `form:"photo_4"`
	Photo5      *multipart.FileHeader `form:"photo_5"`
	Video       *multipart.FileHeader `form:"video"`
}

type CreateRefundResponse struct {
	ID                    uuid.UUID `json:"id"`
	OrderID               uuid.UUID `json:"order_id"`
	Reason                string    `json:"reason"`
	RefundRequestDatetime time.Time `json:"refund_request_datetime"`
	Status                string    `json:"status"`
}

func NewCreateRefundResponse(source *models.Refund) (target *CreateRefundResponse) {
	target = &CreateRefundResponse{
		ID:                    source.ID,
		OrderID:               source.OrderID,
		Reason:                source.Reason,
		RefundRequestDatetime: source.RefundRequestDatetime,
		Status:                source.Status,
	}
	return
}

type GetAllRefundQuery struct {
	SearchInvoice   string     `json:"search-invoice"`
	FilterStatus    string     `json:"filter-status"`
	FilterStartDate *time.Time `json:"filter-start-date"`
	FilterEndDate   *time.Time `json:"filter-end-date"`
	Sort            string     `json:"sort"`
	SortBy          string     `json:"sort-by"`
	Limit           uint       `json:"limit"`
	Page            uint       `json:"page"`
}

func NewGetAllRefundQuery(c *gin.Context) *GetAllRefundQuery {
	var query GetAllRefundQuery

	searchInvoice := c.Query("search-invoice")
	filterStatus := strings.ToLower(c.Query("filter-status"))
	filterStartDate := c.Query("filter-start-date")
	filterEndDate := c.Query("filter-end-date")
	sortBy := strings.ToLower(c.Query("sort-by"))
	sort := strings.ToLower(c.Query("sort"))
	limit := c.Query("limit")
	page := c.Query("page")

	if sortBy == "refund_request_datetime" {
		query.SortBy = "refund_request_datetime"
	} else {
		query.SortBy = "refund_request_datetime"
	}

	if sort == "asc" {
		query.Sort = "ASC"
	} else {
		query.Sort = "DESC"
	}

	query.SearchInvoice = searchInvoice
	query.FilterStatus = filterStatus

	if filterStartDate != "" && filterEndDate != "" {
		startDate, err := time.Parse("2006-01-02", filterStartDate)

		if err != nil {
			query.FilterStartDate = nil
		} else {
			query.FilterStartDate = &startDate
		}

		endDate, err := time.Parse("2006-01-02", filterEndDate)

		if err != nil {
			query.FilterEndDate = nil
		} else {
			query.FilterEndDate = &endDate
		}

	} else {
		query.FilterStartDate = nil
		query.FilterEndDate = nil
	}

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		query.Page = 1
	} else if pageInt <= 0 {
		query.Page = 1
	} else {
		query.Page = uint(pageInt)
	}

	limitInt, err := strconv.Atoi(limit)

	if err != nil {
		query.Limit = 10
	} else if limitInt <= 0 {
		query.Limit = 10
	} else {
		query.Limit = uint(limitInt)
	}

	return &query
}

type GetRefundDetailsResponse struct {
	RefundID                uuid.UUID           `json:"refund_id"`
	OrderID                 uuid.UUID           `json:"order_id"`
	InvoiceNumber           string              `json:"invoice_number"`
	MerchantID              uint                `json:"merchant_id"`
	ShopName                string              `json:"shop_name"`
	UserID                  uuid.UUID           `json:"user_id"`
	UserFullName            string              `json:"user_full_name"`
	RefundRequestDatetime   time.Time           `json:"refund_request_datetime"`
	RejectedDatetime        *time.Time          `json:"rejected_datetime"`
	RefundStatus            string              `json:"refund_status"`
	Images                  pq.StringArray      `json:"images" gorm:"type:text[]"`
	Videos                  pq.StringArray      `json:"videos" gorm:"type:text[]"`
	RefundItems             []*dtos.RefundItems `json:"refund_items"`
	TotalUserRefundAmount   float64             `json:"total_user_refund_amount"`
	TotalMerchantLostAmount float64             `json:"total_merchant_lost_amount"`
}

func NewGetAllRefundDetailsResponse(source []*dtos.RefundDetails) (target []*GetRefundDetailsResponse) {
	for _, result := range source {
		target = append(target, &GetRefundDetailsResponse{
			RefundID:                result.RefundID,
			OrderID:                 result.OrderID,
			InvoiceNumber:           result.InvoiceNumber,
			MerchantID:              result.MerchantID,
			ShopName:                result.ShopName,
			UserID:                  result.UserID,
			UserFullName:            result.UserFullName,
			RefundRequestDatetime:   result.RefundRequestDatetime,
			RejectedDatetime:        result.RejectedDatetime,
			RefundStatus:            result.RefundStatus,
			Images:                  result.Images,
			Videos:                  result.Videos,
			RefundItems:             result.RefundItems,
			TotalUserRefundAmount:   result.TotalUserRefundAmount,
			TotalMerchantLostAmount: result.TotalMerchantLostAmount,
		})
	}

	return
}

func NewGetRefundDetailsResponse(source *dtos.RefundDetails) (target *GetRefundDetailsResponse) {
	target = &GetRefundDetailsResponse{
		RefundID:                source.RefundID,
		OrderID:                 source.OrderID,
		InvoiceNumber:           source.InvoiceNumber,
		MerchantID:              source.MerchantID,
		ShopName:                source.ShopName,
		UserID:                  source.UserID,
		UserFullName:            source.UserFullName,
		RefundRequestDatetime:   source.RefundRequestDatetime,
		RejectedDatetime:        source.RejectedDatetime,
		RefundStatus:            source.RefundStatus,
		Images:                  source.Images,
		Videos:                  source.Videos,
		RefundItems:             source.RefundItems,
		TotalUserRefundAmount:   source.TotalUserRefundAmount,
		TotalMerchantLostAmount: source.TotalMerchantLostAmount,
	}

	return
}

type ApproveRefundResponse struct {
	ID                    uuid.UUID  `json:"id"`
	OrderID               uuid.UUID  `json:"order_id"`
	Reason                string     `json:"reason"`
	RefundRequestDatetime time.Time  `json:"refund_request_datetime"`
	Status                string     `json:"status"`
	RejectedDatetime      *time.Time `json:"rejected_datetime"`
}

func NewApproveRefundResponse(source *models.Refund) *ApproveRefundResponse {
	return &ApproveRefundResponse{
		ID:                    source.ID,
		OrderID:               source.OrderID,
		Reason:                source.Reason,
		RefundRequestDatetime: source.RefundRequestDatetime,
		Status:                source.Status,
		RejectedDatetime:      source.RejectedDatetime,
	}
}
