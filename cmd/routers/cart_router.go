package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
)

func cart(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/carts"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("", func(c *gin.Context) { middlewares.Role(c, "001") }, h.CartGetAllByUserId)
		privateRoute.POST("", func(c *gin.Context) { middlewares.Role(c, "001") }, h.CartUpdate)
		privateRoute.POST("/payment/wallet",
			middlewares.DBTransactionMiddleware(),
			func(c *gin.Context) { middlewares.Role(c, "001") },
			h.PayCartWithWallet)
		privateRoute.POST("/payment/sealabs-pay",
			middlewares.DBTransactionMiddleware(),
			func(c *gin.Context) { middlewares.Role(c, "001") },
			h.PayCartWithSealabsPay)
	}
	privateRoute.POST("/payment/sealabs-pay/success",
		middlewares.DBTransactionMiddleware(),
		func(c *gin.Context) { middlewares.Role(c, "001") },
		h.PayCartWithSealabsPaySuccess)
}
