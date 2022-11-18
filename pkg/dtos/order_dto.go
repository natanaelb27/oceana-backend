package dtos

import (
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type OrderListQuery struct {
	FilterOrderStatus          string     `json:"filter_order_status"`
	FilterOngoingOrderStatus   string     `json:"filter_ongoing_order_status"`
	FilterStartTransactionDate *time.Time `json:"filter_start_transaction_date"`
	FilterEndTransactionDate   *time.Time `json:"filter_end_transaction_date"`
	Limit                      uint       `json:"limit"`
	Page                       uint       `json:"page"`
}

func NewOrderListQuery(c *gin.Context) *OrderListQuery {
	var query OrderListQuery

	filterOrderStatus := c.Query("filter_order_status")
	filterOngoingOrderStatus := c.Query("filter_ongoing_order_status")
	filterStartTransactionDate := c.Query("filter_start_transaction_date")
	filterEndTransactionDate := c.Query("filter_end_transaction_date")

	limit := c.Query("limit")
	page := c.Query("page")

	query.FilterOrderStatus = filterOrderStatus
	query.FilterOngoingOrderStatus = filterOngoingOrderStatus
	startDate, err := time.Parse("2006-01-02", filterStartTransactionDate)

	if err != nil {
		query.FilterStartTransactionDate = nil
	} else {
		query.FilterStartTransactionDate = &startDate
	}

	endDate, err := time.Parse("2006-01-02", filterEndTransactionDate)

	if err != nil {
		query.FilterEndTransactionDate = nil
	} else {
		query.FilterEndTransactionDate = &endDate
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

type BuyerOrderListGetAllResponse struct {
	Pagination *utils.Pagination `json:"pagination"`
	OrderList  []*dtos.OrderList `json:"order_list"`
}

func NewBuyerOrderListGetAllResponse(source []*dtos.OrderList, pagination *utils.Pagination) (target *BuyerOrderListGetAllResponse) {
	if pagination == nil {
		return nil
	}

	target = &BuyerOrderListGetAllResponse{
		Pagination: pagination,
		OrderList:  source,
	}

	return
}
