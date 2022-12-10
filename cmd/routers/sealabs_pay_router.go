package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func sealabsPay(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	externalPath := "/external/sealabs-pay"
	externalRoute := public.Group(externalPath)
	{
		externalRoute.POST("/top-up", h.SealabsPayTopupCallback)
	}

	internalPath := "/sealabs-pay"
	publicRoute := public.Group(internalPath)
	{
		publicRoute.POST("/pay", h.SealabsPayTopup)
	}
}
