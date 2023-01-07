package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
)

func admin(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/admin"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("/refunds", func(c *gin.Context) { middlewares.Role(c, "100") }, h.GetAllRefundDetailsByAdmin)
		privateRoute.PATCH("/refunds/:refund-id/complete-refund", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "100") }, h.UpdateCompleteRefund)
	}
}
