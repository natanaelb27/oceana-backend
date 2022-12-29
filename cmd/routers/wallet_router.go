package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
)

func wallet(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/wallet"
	privateRoute := private.Group(path)
	{
		privateRoute.POST("/topup", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "001") }, h.WalletTopup)
	}
}
