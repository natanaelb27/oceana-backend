package dtos

import (
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	if filterStartTransactionDate != "" && filterEndTransactionDate != "" {
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

	} else {
		query.FilterStartTransactionDate = nil
		query.FilterEndTransactionDate = nil
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

type OrderListGetAllResponse struct {
	Pagination *utils.Pagination `json:"pagination"`
	OrderList  []*dtos.OrderList `json:"order_list"`
}

func NewOrderListGetAllResponse(source []*dtos.OrderList, pagination *utils.Pagination) (target *OrderListGetAllResponse) {
	if pagination == nil {
		return nil
	}

	target = &OrderListGetAllResponse{
		Pagination: pagination,
		OrderList:  source,
	}

	return
}

type OrderListDetailsGetAllResponse struct {
	TransactionID                     uuid.UUID                    `json:"transaction_id"`
	OrderID                           uuid.UUID                    `json:"oder_id"`
	MerchantID                        uint                         `json:"merchant_id"`
	ShopName                          string                       `json:"shop_name"`
	OrderDetails                      []*dtos.OrderDetailList      `json:"order_details"`
	TransactionDatetime               time.Time                    `json:"transaction_datetime"`
	OrderStatus                       string                       `json:"order_status"`
	CourierDescription                string                       `json:"courier_description"`
	Address                           *dtos.AddressSnapshotDetails `json:"address"`
	PaymentMethod                     string                       `json:"payment_method"`
	TotalPriceBeforeDiscount          uint                         `json:"total_price_before_discount"`
	TotalMerchantPromotionDiscount    float64                      `json:"total_merchant_promotion_discount"`
	TotalMerchantVoucherDiscount      float64                      `json:"total_merchant_voucher_discount"`
	TotalMarketplacePromotionDiscount float64                      `json:"total_marketplace_promotion_discount"`
	TotalMarketplaceVoucherDiscount   float64                      `json:"total_marketplace_voucher_discount"`
	ShippingFee                       uint                         `json:"shipping_fee"`
	TotalPayment                      float64                      `json:"total_payment"`
}

func NewOrderListDetailsGetAllResponse(source *dtos.OrderListDetails) (target *OrderListDetailsGetAllResponse) {
	target = &OrderListDetailsGetAllResponse{
		TransactionID:                     source.TransactionID,
		OrderID:                           source.OrderID,
		MerchantID:                        source.MerchantID,
		ShopName:                          source.ShopName,
		OrderDetails:                      source.OrderDetails,
		TransactionDatetime:               source.TransactionDatetime,
		OrderStatus:                       source.OrderStatus,
		CourierDescription:                source.CourierDescription,
		Address:                           source.Address,
		PaymentMethod:                     source.PaymentMethod,
		TotalPriceBeforeDiscount:          source.TotalPriceBeforeDiscount,
		TotalMerchantPromotionDiscount:    source.TotalMerchantPromotionDiscount,
		TotalMerchantVoucherDiscount:      source.TotalMerchantVoucherDiscount,
		TotalMarketplacePromotionDiscount: source.TotalMarketplacePromotionDiscount,
		TotalMarketplaceVoucherDiscount:   source.TotalMarketplaceVoucherDiscount,
		ShippingFee:                       source.ShippingFee,
		TotalPayment:                      source.TotalPayment,
	}

	return
}
