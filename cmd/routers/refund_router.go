package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
)

func refund(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/refund"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("/:refund-id", func(c *gin.Context) { middlewares.Role(c, "111") }, h.GetDetailsByID)
	}
}
