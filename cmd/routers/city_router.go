package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func city(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/cities"
	publicRoute := public.Group(path)
	{
		publicRoute.GET("", h.CitiesGetAll)

	}
}
