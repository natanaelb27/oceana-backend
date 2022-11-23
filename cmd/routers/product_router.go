package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func product(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/products"
	publicRoute := public.Group(path)
	{
		publicRoute.GET("", h.ProductGetAll)
	}
}
