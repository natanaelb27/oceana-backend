package dtos

import (
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
