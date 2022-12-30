package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
)

func admin(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/admin"
	privateRoute := private.Group(path)
	{
		privateRoute.GET("/refunds", func(c *gin.Context) { middlewares.Role(c, "100") }, h.GetAllRefundDetailsByAdmin)
		privateRoute.PATCH("/refunds/:refund-id/complete-refund", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "100") }, h.UpdateCompleteRefund)
	}
}
