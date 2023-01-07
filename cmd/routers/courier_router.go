package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
)

func courier(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/couriers"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("", func(c *gin.Context) { middlewares.Role(c, "011") }, h.CourierGetAll)
	}
}
