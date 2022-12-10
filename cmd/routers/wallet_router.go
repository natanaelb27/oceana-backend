package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func wallet(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/wallet"
	privateRoute := private.Group(path)
	{
		privateRoute.POST("/top-up", h.WalletTopup)
	}
}
