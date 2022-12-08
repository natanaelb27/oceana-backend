package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
)

func cart(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/carts"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("", func(c *gin.Context) { middlewares.Role(c, "001") }, h.CartGetAllByUserId)
		privateRoute.POST("", func(c *gin.Context) { middlewares.Role(c, "001") }, h.CartUpdate)
		privateRoute.POST("/payment/wallet", func(c *gin.Context) { middlewares.Role(c, "001") }, h.PayCartWithWallet)
	}
}
