package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func sealabsPay(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/external/sealabs-pay"
	publicRoute := public.Group(path)
	{
		publicRoute.POST("", h.SealabsPayTest)
	}
}
