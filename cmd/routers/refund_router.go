package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
)

func refund(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/refunds"

	publicRoute := public.Group(path)
	{
		publicRoute.GET("/calculate-refund-details/:order-detail-id", h.CalculateRefundDetails)
	}

	privateRoute := private.Group(path)
	{
		privateRoute.GET("/:refund-id", func(c *gin.Context) { middlewares.Role(c, "111") }, h.GetDetailsByID)
	}
}
