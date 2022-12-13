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
		// vouchers
		privateRoute.POST("/vouchers", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantVoucherCreate)
		privateRoute.PUT("/vouchers/:merchant-voucher-id", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantVoucherUpdate)
		privateRoute.DELETE("/vouchers/:merchant-voucher-id", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantVoucherDelete)
		// promotions
		privateRoute.POST("/promotions", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionCreate)
		privateRoute.GET("/promotions", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionGetAll)
		privateRoute.GET("/promotions/:merchant-promotion-id", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionGetByID)
		privateRoute.PUT("/promotions/:merchant-promotion-id", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionUpdate)
		// orders
		privateRoute.PATCH("/orders/:order-id", func(c *gin.Context) { middlewares.Role(c, "010") }, h.OrderProcessOrder)
		privateRoute.GET("/transaction-history", func(c *gin.Context) { middlewares.Role(c, "010") }, h.OrderListGetAllByMerchantID)
		privateRoute.GET("/transaction-history/:orderId", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantOrderListDetailsGetByID)
		// logistic
		privateRoute.GET("/logistics", func(c *gin.Context) { middlewares.Role(c, "010") }, h.LogisticGet)
		privateRoute.PATCH("/logistics", func(c *gin.Context) { middlewares.Role(c, "010") }, h.LogisticUpdate)
	}
}
