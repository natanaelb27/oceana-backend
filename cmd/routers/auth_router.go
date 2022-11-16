package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func auth(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/auth"
	publicRoute := public.Group(path)
	{
		publicRoute.POST("/register", h.UserCreate)
	}
}
