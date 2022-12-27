package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
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
		privateRoute.POST("/promotions", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionCreate)
		privateRoute.GET("/promotions", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionGetAll)
		privateRoute.GET("/promotions/:merchant-promotion-id", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionGetByID)
		privateRoute.PUT("/promotions/:merchant-promotion-id", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionUpdate)
		privateRoute.DELETE("/promotions/:merchant-promotion-id", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionDelete)
		// orders
		privateRoute.PATCH("/orders/:order-id/process-delivery", func(c *gin.Context) { middlewares.Role(c, "010") }, h.OrderProcessDelivery)
		privateRoute.PATCH("/orders/:order-id/process-order", func(c *gin.Context) { middlewares.Role(c, "010") }, h.OrderProcessOrder)
		privateRoute.PATCH("/orders/:order-id/cancel-order", func(c *gin.Context) { middlewares.Role(c, "010") }, h.SellerCancelOrder)
		privateRoute.PATCH("/orders/:order-id/claim-funds", func(c *gin.Context) { middlewares.Role(c, "010") }, h.SellerClaimFunds)
		privateRoute.GET("/transaction-history", func(c *gin.Context) { middlewares.Role(c, "010") }, h.OrderListGetAllByMerchantID)
		privateRoute.GET("/transaction-history/:orderId", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantOrderListDetailsGetByID)

		// refunds
		privateRoute.PATCH("/refunds/:refund-id/approve", func(c *gin.Context) { middlewares.Role(c, "010") }, h.UpdateApproveRefund)
		privateRoute.PATCH("/refunds/:refund-id/reject", func(c *gin.Context) { middlewares.Role(c, "010") }, h.UpdateRejectRefund)

		// product
		privateRoute.PUT("/products/:product_id",
			middlewares.DBTransactionMiddleware(),
			func(c *gin.Context) { middlewares.Role(c, "011") },
			h.ProductUpdate,
		)
		privateRoute.POST("/products",
			middlewares.DBTransactionMiddleware(),
			func(c *gin.Context) { middlewares.Role(c, "010") },
			h.ProductCreate)
		privateRoute.GET("/products",
			middlewares.DBTransactionMiddleware(),
			func(c *gin.Context) { middlewares.Role(c, "010") },
			h.MyShopGetAll)

		// logistic
		privateRoute.GET("/logistics", func(c *gin.Context) { middlewares.Role(c, "010") }, h.LogisticGet)
		privateRoute.PATCH("/logistics", func(c *gin.Context) { middlewares.Role(c, "010") }, h.LogisticUpdate)
		privateRoute.GET("/logistics/selected", func(c *gin.Context) { middlewares.Role(c, "010") }, h.LogisticGetSelected)

		// refunds
		privateRoute.GET("/refunds", func(c *gin.Context) { middlewares.Role(c, "010") }, h.GetAllRefundDetailsByMerchant)
	}
}
