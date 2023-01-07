package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
)

func city(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/cities"
	publicRoute := public.Group(path)
	{
		publicRoute.GET("", h.CitiesGetAll)

	}
}
