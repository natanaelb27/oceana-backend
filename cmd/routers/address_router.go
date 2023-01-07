package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
)

func address(private *gin.RouterGroup, h *handlers.Handler) {
	path := "/profile/addresses"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("", func(c *gin.Context) { middlewares.Role(c, "001") }, h.AddressGetAll)
		privateRoute.GET("/:address_id", func(c *gin.Context) { middlewares.Role(c, "001") }, h.AddressGetById)
		privateRoute.POST("", func(c *gin.Context) { middlewares.Role(c, "001") }, h.AddressCreate)
		privateRoute.PUT("/:address_id", func(c *gin.Context) { middlewares.Role(c, "001") }, h.AddressUpdate)
		privateRoute.DELETE("/:address_id", func(c *gin.Context) { middlewares.Role(c, "001") }, h.AddressDelete)
		privateRoute.PATCH("/:address_id/default", func(c *gin.Context) { middlewares.Role(c, "001") }, h.AddressUpdateDefault)
		privateRoute.PATCH("/:address_id/default-shop", func(c *gin.Context) { middlewares.Role(c, "011") }, h.AddressShopUpdateDefault)
	}
}
