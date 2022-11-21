package dtos

// type ProductDetailsGetByIDResponse struct {
// 	TransactionID                   uuid.UUID                    `json:"transaction_id"`
// 	OrderID                         uuid.UUID                    `json:"oder_id"`
// 	MerchantID                      uint                         `json:"merchant_id"`
// 	ShopName                        string                       `json:"shop_name"`
// 	OrderDetails                    []*dtos.OrderDetailList      `json:"order_details"`
// 	TransactionDatetime             time.Time                    `json:"transaction_datetime"`
// 	OrderStatus                     string                       `json:"order_status"`
// 	CourierDescription              string                       `json:"courier_description"`
// 	OrderTrackingHistory            []*dtos.OrderTrackingHistory `json:"order_tracking_history" gorm:"-"`
// 	Address                         *dtos.AddressSnapshotDetails `json:"address"`
// 	PaymentMethod                   string                       `json:"payment_method"`
// 	TotalPriceBeforeDiscount        uint                         `json:"total_price_before_discount"`
// 	TotalMerchantPromotionDiscount  float64                      `json:"total_merchant_promotion_discount"`
// 	TotalMerchantVoucherDiscount    float64                      `json:"total_merchant_voucher_discount"`
// 	TotalMarketplaceVoucherDiscount float64                      `json:"total_marketplace_voucher_discount"`
// 	ShippingFee                     uint                         `json:"shipping_fee"`
// 	TotalPayment                    float64                      `json:"total_payment"`
// }

// func NewProductDetailsGetByIDResponse(source *dtos.OrderListDetails) (target *OrderListDetailsGetByIDResponse) {
// 	target = &OrderListDetailsGetByIDResponse{
// 		TransactionID:                   source.TransactionID,
// 		OrderID:                         source.OrderID,
// 		MerchantID:                      source.MerchantID,
// 		ShopName:                        source.ShopName,
// 		OrderDetails:                    source.OrderDetails,
// 		TransactionDatetime:             source.TransactionDatetime,
// 		OrderStatus:                     source.OrderStatus,
// 		CourierDescription:              source.CourierDescription,
// 		OrderTrackingHistory:            source.OrderTrackingHistory,
// 		Address:                         source.Address,
// 		PaymentMethod:                   source.PaymentMethod,
// 		TotalPriceBeforeDiscount:        source.TotalPriceBeforeDiscount,
// 		TotalMerchantPromotionDiscount:  source.TotalMerchantPromotionDiscount,
// 		TotalMerchantVoucherDiscount:    source.TotalMerchantVoucherDiscount,
// 		TotalMarketplaceVoucherDiscount: source.TotalMarketplaceVoucherDiscount,
// 		ShippingFee:                     source.ShippingFee,
// 		TotalPayment:                    source.TotalPayment,
// 	}

// 	return
// }
