package dtos

import (
	"strconv"
	"strings"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderListQuery struct {
	Search                     string     `json:"string"`
	FilterOrderStatus          string     `json:"filter-order-status"`
	FilterOngoingOrderStatus   string     `json:"filter-ongoing-order-status"`
	FilterStartTransactionDate *time.Time `json:"filter-start-transaction-date"`
	FilterEndTransactionDate   *time.Time `json:"filter-end-transaction-date"`
	FilterMerchantID           uint       `json:"filter-merchant-id"`
	Limit                      uint       `json:"limit"`
	Page                       uint       `json:"page"`
}

func NewOrderListQuery(c *gin.Context) *OrderListQuery {
	var query OrderListQuery

	search := c.Query("search")
	filterOrderStatus := strings.ToLower(c.Query("filter-order-status"))
	filterOngoingOrderStatus := strings.ToLower(c.Query("filter-ongoing-order-status"))
	filterStartTransactionDate := c.Query("filter-start-transaction-date")
	filterEndTransactionDate := c.Query("filter-end-transaction-date")
	filterMerchantID := utils.StringToUint(c.Query("filter-merchant-id"))

	limit := c.Query("limit")
	page := c.Query("page")

	query.Search = search
	query.FilterOrderStatus = filterOrderStatus
	query.FilterOngoingOrderStatus = filterOngoingOrderStatus
	query.FilterMerchantID = filterMerchantID

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
	OrderList []*dtos.OrderList `json:"order_list"`
}

func NewOrderListGetAllResponse(source []*dtos.OrderList) (target *OrderListGetAllResponse) {

	target = &OrderListGetAllResponse{
		OrderList: source,
	}

	return
}

type OrderListDetailsGetByIDResponse struct {
	TransactionID                   uuid.UUID                    `json:"transaction_id"`
	OrderID                         uuid.UUID                    `json:"oder_id"`
	InvoiceNumber                   string                       `json:"invoice_number"`
	MerchantID                      uint                         `json:"merchant_id"`
	ShopName                        string                       `json:"shop_name"`
	OrderDetails                    []*dtos.OrderDetailList      `json:"order_details"`
	TransactionDatetime             time.Time                    `json:"transaction_datetime"`
	OrderStatus                     string                       `json:"order_status"`
	CourierDescription              string                       `json:"courier_description"`
	OrderTrackingHistory            []*dtos.OrderTrackingHistory `json:"order_tracking_history" gorm:"-"`
	Address                         *dtos.AddressSnapshotDetails `json:"address"`
	PaymentMethod                   string                       `json:"payment_method"`
	TotalPriceBeforeDiscount        uint                         `json:"total_price_before_discount"`
	TotalMerchantPromotionDiscount  float64                      `json:"total_merchant_promotion_discount"`
	TotalMerchantVoucherDiscount    float64                      `json:"total_merchant_voucher_discount"`
	TotalMarketplaceVoucherDiscount float64                      `json:"total_marketplace_voucher_discount"`
	ShippingFee                     uint                         `json:"shipping_fee"`
	TotalPayment                    float64                      `json:"total_payment"`
}

func NewOrderListDetailsGetByIDResponse(source *dtos.OrderListDetails) (target *OrderListDetailsGetByIDResponse) {
	target = &OrderListDetailsGetByIDResponse{
		TransactionID:                   source.TransactionID,
		OrderID:                         source.OrderID,
		InvoiceNumber:                   source.InvoiceNumber,
		MerchantID:                      source.MerchantID,
		ShopName:                        source.ShopName,
		OrderDetails:                    source.OrderDetails,
		TransactionDatetime:             source.TransactionDatetime,
		OrderStatus:                     source.OrderStatus,
		CourierDescription:              source.CourierDescription,
		OrderTrackingHistory:            source.OrderTrackingHistory,
		Address:                         source.Address,
		PaymentMethod:                   source.PaymentMethod,
		TotalPriceBeforeDiscount:        source.TotalPriceBeforeDiscount,
		TotalMerchantPromotionDiscount:  source.TotalMerchantPromotionDiscount,
		TotalMerchantVoucherDiscount:    source.TotalMerchantVoucherDiscount,
		TotalMarketplaceVoucherDiscount: source.TotalMarketplaceVoucherDiscount,
		ShippingFee:                     source.ShippingFee,
		TotalPayment:                    source.TotalPayment,
	}

	return
}

type OrderProcessOrderRequest struct {
	CourierServiceID uint `json:"courier_service_id" binding:"required"`
}

type OrderProcessOrderResponse struct {
	ID               uuid.UUID `json:"id"`
	CourierServiceID uint      `json:"courier_service_id"`
	Status           string    `json:"status"`
}

func NewOrderProcessOrderResponse(source *models.Order) *OrderProcessOrderResponse {
	return &OrderProcessOrderResponse{
		ID:               source.ID,
		CourierServiceID: source.CourierServiceID,
		Status:           source.Status,
	}
}
