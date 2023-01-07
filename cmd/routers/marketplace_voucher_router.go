package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
)

func marketplaceVoucher(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/marketplace/vouchers"
	publicRoute := public.Group(path)
	{
		publicRoute.GET("/active", h.MarketplaceVoucherGetAllByActive)
	}
	privateRoute := private.Group(path)
	{
		privateRoute.GET("", func(c *gin.Context) { middlewares.Role(c, "100") }, h.MarketplaceVoucherGetAll)
		privateRoute.GET("/:id", func(c *gin.Context) { middlewares.Role(c, "100") }, h.MarketplaceVoucherGetById)
		privateRoute.POST("", func(c *gin.Context) { middlewares.Role(c, "100") }, h.MarketplaceVoucherCreate)
		privateRoute.DELETE("/:id", func(c *gin.Context) { middlewares.Role(c, "100") }, h.MarketplaceVoucherDelete)
		privateRoute.PUT("/:id", func(c *gin.Context) { middlewares.Role(c, "100") }, h.MarketplaceVoucherUpdate)
		privateRoute.GET("/buyer", func(c *gin.Context) { middlewares.Role(c, "001") }, h.MarketplaceVoucherGetAllByBuyer)
	}
}
