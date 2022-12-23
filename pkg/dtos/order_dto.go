package dtos

import (
	"strconv"
	"strings"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderListQuery struct {
	Search                     string     `json:"string"`
	FilterOrderStatus          string     `json:"filter-order-status"`
	FilterOngoingOrderStatus   string     `json:"filter-ongoing-order-status"`
	FilterStartTransactionDate *time.Time `json:"filter-start-transaction-date"`
	FilterEndTransactionDate   *time.Time `json:"filter-end-transaction-date"`
	FilterIsReviewed           string     `json:"filter-is-reviewed"`
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
	filterIsReviewed := c.Query("filter-is-reviewed")

	limit := c.Query("limit")
	page := c.Query("page")

	query.Search = search
	query.FilterOrderStatus = filterOrderStatus
	query.FilterOngoingOrderStatus = filterOngoingOrderStatus

	if filterIsReviewed != "" {
		if filterIsReviewed != "true" && filterIsReviewed != "false" {
			query.FilterIsReviewed = "all"
		} else {
			query.FilterIsReviewed = filterIsReviewed
		}
	}

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
	TransactionID                        uuid.UUID               `json:"transaction_id"`
	OrderID                              uuid.UUID               `json:"order_id"`
	InvoiceNumber                        string                  `json:"invoice_number"`
	MerchantID                           uint                    `json:"merchant_id"`
	ShopName                             string                  `json:"shop_name"`
	ShopPhoto                            string                  `json:"shop_photo"`
	UserID                               uuid.UUID               `json:"user_id"`
	UserFullName                         string                  `json:"user_full_name"`
	CourierServiceID                     uint                    `json:"courier_service_id"`
	CourierDescription                   string                  `json:"courier_description"`
	OrderDetails                         []*dtos.OrderDetailList `json:"order_details"`
	TransactionDatetime                  time.Time               `json:"transaction_datetime"`
	OrderStatus                          string                  `json:"order_status"`
	TotalPrice                           float64                 `json:"total_price"`
	RefundID                             *uuid.UUID              `json:"refund_id"`
	ClaimStatus                          bool                    `json:"claim_status"`
	IsAllReviewed                        bool                    `json:"is_all_reviewed"`
	SellerConfirmationExpirationDatetime *time.Time              `json:"seller_confirmation_expiration_datetime"`
	OnProcessExpirationDatetime          *time.Time              `json:"on_process_expiration_datetime"`
	DeliveredExpirationDatetime          *time.Time              `json:"delivered_expiration_datetime"`
}

func NewOrderListGetAllResponse(source []*dtos.OrderList) (target []*OrderListGetAllResponse) {
	for _, result := range source {
		target = append(target, &OrderListGetAllResponse{
			TransactionID:                        result.TransactionID,
			OrderID:                              result.OrderID,
			InvoiceNumber:                        result.InvoiceNumber,
			MerchantID:                           result.MerchantID,
			ShopName:                             result.ShopName,
			ShopPhoto:                            result.ShopPhoto,
			UserID:                               result.UserID,
			UserFullName:                         result.UserFullName,
			CourierServiceID:                     result.CourierServiceID,
			CourierDescription:                   result.CourierDescription,
			OrderDetails:                         result.OrderDetails,
			TransactionDatetime:                  result.TransactionDatetime,
			OrderStatus:                          result.OrderStatus,
			TotalPrice:                           result.TotalPrice,
			RefundID:                             result.RefundID,
			ClaimStatus:                          result.ClaimStatus,
			IsAllReviewed:                        result.IsAllReviewed,
			SellerConfirmationExpirationDatetime: result.SellerConfirmationExpirationDatetime,
			OnProcessExpirationDatetime:          result.OnProcessExpirationDatetime,
			DeliveredExpirationDatetime:          result.DeliveredExpirationDatetime,
		})
	}

	return
}

type OrderListDetailsGetByIDResponse struct {
	TransactionID                        uuid.UUID                    `json:"transaction_id"`
	OrderID                              uuid.UUID                    `json:"oder_id"`
	InvoiceNumber                        string                       `json:"invoice_number"`
	MerchantID                           uint                         `json:"merchant_id"`
	ShopName                             string                       `json:"shop_name"`
	ShopPhoto                            string                       `json:"shop_photo"`
	UserID                               uuid.UUID                    `json:"user_id"`
	UserFullName                         string                       `json:"user_full_name"`
	OrderDetails                         []*dtos.OrderDetailList      `json:"order_details"`
	TransactionDatetime                  time.Time                    `json:"transaction_datetime"`
	OrderStatus                          string                       `json:"order_status"`
	CourierServiceID                     uint                         `json:"courier_service_id"`
	CourierDescription                   string                       `json:"courier_description"`
	OrderTrackingHistory                 []*dtos.OrderTrackingHistory `json:"order_tracking_history" gorm:"-"`
	Address                              *dtos.AddressSnapshotDetails `json:"address"`
	PaymentMethod                        string                       `json:"payment_method"`
	TotalPriceBeforeDiscount             uint                         `json:"total_price_before_discount"`
	TotalMerchantPromotionDiscount       float64                      `json:"total_merchant_promotion_discount"`
	TotalMerchantVoucherDiscount         float64                      `json:"total_merchant_voucher_discount"`
	TotalMarketplaceVoucherDiscount      float64                      `json:"total_marketplace_voucher_discount"`
	ShippingFee                          uint                         `json:"shipping_fee"`
	TotalPayment                         float64                      `json:"total_payment"`
	RefundID                             *uuid.UUID                   `json:"refund_id"`
	ClaimStatus                          bool                         `json:"claim_status"`
	IsAllReviewed                        bool                         `json:"is_all_reviewed"`
	SellerConfirmationExpirationDatetime *time.Time                   `json:"seller_confirmation_expiration_datetime"`
	OnProcessExpirationDatetime          *time.Time                   `json:"on_process_expiration_datetime"`
	DeliveredExpirationDatetime          *time.Time                   `json:"delivered_expiration_datetime"`
}

func NewOrderListDetailsGetByIDResponse(source *dtos.OrderListDetails) (target *OrderListDetailsGetByIDResponse) {
	target = &OrderListDetailsGetByIDResponse{
		TransactionID:                        source.TransactionID,
		OrderID:                              source.OrderID,
		InvoiceNumber:                        source.InvoiceNumber,
		MerchantID:                           source.MerchantID,
		ShopName:                             source.ShopName,
		ShopPhoto:                            source.ShopPhoto,
		UserID:                               source.UserID,
		UserFullName:                         source.UserFullName,
		OrderDetails:                         source.OrderDetails,
		TransactionDatetime:                  source.TransactionDatetime,
		OrderStatus:                          source.OrderStatus,
		CourierServiceID:                     source.CourierServiceID,
		CourierDescription:                   source.CourierDescription,
		OrderTrackingHistory:                 source.OrderTrackingHistory,
		Address:                              source.Address,
		PaymentMethod:                        source.PaymentMethod,
		TotalPriceBeforeDiscount:             source.TotalPriceBeforeDiscount,
		TotalMerchantPromotionDiscount:       source.TotalMerchantPromotionDiscount,
		TotalMerchantVoucherDiscount:         source.TotalMerchantVoucherDiscount,
		TotalMarketplaceVoucherDiscount:      source.TotalMarketplaceVoucherDiscount,
		ShippingFee:                          source.ShippingFee,
		TotalPayment:                         source.TotalPayment,
		RefundID:                             source.RefundID,
		ClaimStatus:                          source.ClaimStatus,
		IsAllReviewed:                        source.IsAllReviewed,
		SellerConfirmationExpirationDatetime: source.SellerConfirmationExpirationDatetime,
		OnProcessExpirationDatetime:          source.OnProcessExpirationDatetime,
		DeliveredExpirationDatetime:          source.DeliveredExpirationDatetime,
	}

	return
}

type OrderProcessDeliveryResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func NewOrderProcessDeliveryResponse(source *models.Order) *OrderProcessDeliveryResponse {
	return &OrderProcessDeliveryResponse{
		ID:     source.ID,
		Status: source.Status,
	}
}

type UpdateOrderResponse struct {
	ID                      uuid.UUID `json:"id"`
	InvoiceNumber           string    `json:"invoice_number"`
	TransactionID           uuid.UUID `json:"transaction_id"`
	MerchantID              uint      `json:"merchant_id"`
	AddressSnapshotID       uint      `json:"address_snapshot_id"`
	CourierServiceID        uint      `json:"courier_service_id"`
	TotalPrice              float64   `json:"total_price"`
	TotalPriceAfterDiscount float64   `json:"total_price_after_discount"`
	ShippingFee             uint      `json:"shipping_fee"`
	Status                  string    `json:"status"`
	ClaimStatus             bool      `json:"claim_status"`
}

func NewUpdateOrderResponse(source *models.Order) *UpdateOrderResponse {
	return &UpdateOrderResponse{
		ID:                      source.ID,
		InvoiceNumber:           source.InvoiceNumber,
		TransactionID:           source.TransactionID,
		MerchantID:              source.MerchantID,
		AddressSnapshotID:       source.AddressSnapshotID,
		CourierServiceID:        source.CourierServiceID,
		TotalPrice:              source.TotalPrice,
		TotalPriceAfterDiscount: source.TotalPriceAfterDiscount,
		ShippingFee:             source.ShippingFee,
		Status:                  source.Status,
		ClaimStatus:             source.ClaimStatus,
	}
}
