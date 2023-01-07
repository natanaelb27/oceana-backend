package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
)

func wallet(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/wallet"
	privateRoute := private.Group(path)
	{
		privateRoute.POST("/topup", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "001") }, h.WalletTopup)
	}
}
