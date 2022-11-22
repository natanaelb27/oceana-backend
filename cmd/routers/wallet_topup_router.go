package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func walletTopup(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/wallet-top-up"
	publicRoute := public.Group(path)
	{
		publicRoute.POST("/callback", h.WalletTopUpCallback)
	}
}
