package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func sealabsPay(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	callbackPath := "/callback/sealabs-pay"
	callbackRoute := public.Group(callbackPath)
	{
		callbackRoute.POST("/topup", h.SealabsPayTopupCallback)
	}

	path := "/sealabs-pay"
	privateRoute := private.Group(path)
	{
		privateRoute.POST("/topup", h.SealabsPayTopup)
	}
}
