package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
)

func subDistrict(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/sub-districts"
	publicRoute := public.Group(path)
	{
		publicRoute.GET("", h.SubDistrictGetAll)

	}
}
