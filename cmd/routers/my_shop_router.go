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
		privateRoute.POST("/vouchers", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantVoucherCreate)
		privateRoute.PUT("/vouchers/:merchantVoucherId", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantVoucherUpdate)
		privateRoute.DELETE("/vouchers/:merchantVoucherId", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantVoucherDelete)
		privateRoute.POST("/promotions", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionCreate)
		privateRoute.GET("/promotions", func(c *gin.Context) { middlewares.Role(c, "010") }, h.MerchantPromotionGetAll)
	}
}
