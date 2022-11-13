package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func role(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/roles"
	publicRoute := public.Group(path)
	{
		publicRoute.GET("/", h.RoleGetAll)
		publicRoute.GET("/:id", h.RoleGetById)
		publicRoute.POST("", h.RoleCreate)
		publicRoute.DELETE("/:id", h.RoleDelete)
		publicRoute.PUT("/:id", h.RoleUpdate)
	}
	privateRoute := private.Group(path)
	{
		privateRoute.GET("/private-example", h.RoleGetAll)
	}
}
