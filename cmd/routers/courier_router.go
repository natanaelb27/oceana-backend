package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
)

func courier(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/couriers"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("", func(c *gin.Context) { middlewares.Role(c, "011") }, h.CourierGetAll)
	}
}
