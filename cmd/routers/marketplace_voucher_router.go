package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
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
	}
}
