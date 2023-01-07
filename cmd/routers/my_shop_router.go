package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
)

func myShop(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/my-shop"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("/profile", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantGetProfile)
		privateRoute.PATCH("/profile/photo", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantUpdatePhoto)
		// vouchers
		privateRoute.POST("/vouchers", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantVoucherCreate)
		privateRoute.PUT("/vouchers/:merchant-voucher-id", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantVoucherUpdate)
		privateRoute.DELETE("/vouchers/:merchant-voucher-id", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantVoucherDelete)
		// promotions
		privateRoute.POST("/promotions",
			middlewares.DBTransactionMiddleware(),
			func(c *gin.Context) { middlewares.Role(c, "010") },
			h.MerchantPromotionCreate)
		privateRoute.GET("/promotions", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionGetAll)
		privateRoute.GET("/promotions/:merchant-promotion-id", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionGetByID)
		privateRoute.PUT("/promotions/:merchant-promotion-id",
			middlewares.DBTransactionMiddleware(),
			func(c *gin.Context) { middlewares.Role(c, "010") },
			h.MerchantPromotionUpdate)
		privateRoute.DELETE("/promotions/:merchant-promotion-id", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionDelete)
		// orders
		privateRoute.PATCH("/orders/:order-id/process-delivery", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "010") }, h.OrderProcessDelivery)
		privateRoute.PATCH("/orders/:order-id/process-order", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "010") }, h.OrderProcessOrder)
		privateRoute.PATCH("/orders/:order-id/cancel-order", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "010") }, h.SellerCancelOrder)
		privateRoute.PATCH("/orders/:order-id/claim-funds", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "010") }, h.SellerClaimFunds)
		privateRoute.GET("/transaction-history", func(c *gin.Context) { middlewares.Role(c, "010") }, h.OrderListGetAllByMerchantID)
		privateRoute.GET("/transaction-history/:orderId", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantOrderListDetailsGetByID)

		// refunds
		privateRoute.PATCH("/refunds/:refund-id/approve", func(c *gin.Context) { middlewares.Role(c, "010") }, h.UpdateApproveRefund)
		privateRoute.PATCH("/refunds/:refund-id/reject", func(c *gin.Context) { middlewares.Role(c, "010") }, h.UpdateRejectRefund)

		// product
		privateRoute.PATCH("/products/:product_id",
			func(c *gin.Context) { middlewares.Role(c, "011") },
			middlewares.DBTransactionMiddleware(),
			h.ProductUpdate,
		)
		privateRoute.POST("/products",
			func(c *gin.Context) { middlewares.Role(c, "010") },
			middlewares.DBTransactionMiddleware(),
			h.ProductCreate)
		privateRoute.GET("/products",
			func(c *gin.Context) { middlewares.Role(c, "010") },
			middlewares.DBTransactionMiddleware(),
			h.MyShopGetAll)

		// logistic
		privateRoute.GET("/logistics", func(c *gin.Context) { middlewares.Role(c, "010") }, h.LogisticGet)
		privateRoute.PATCH("/logistics", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "010") }, h.LogisticUpdate)
		privateRoute.GET("/logistics/selected", func(c *gin.Context) { middlewares.Role(c, "010") }, h.LogisticGetSelected)

		// refunds
		privateRoute.GET("/refunds", func(c *gin.Context) { middlewares.Role(c, "010") }, h.GetAllRefundDetailsByMerchant)
	}
}
