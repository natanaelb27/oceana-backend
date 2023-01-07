package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
)

func review(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/reviews"
	publicRoute := public.Group(path)
	{
		publicRoute.GET("products/:productID", h.ProductReviewGet)
	}
	privateRoute := private.Group(path)
	{
		privateRoute.POST("/order-details/:orderDetailID", middlewares.DBTransactionMiddleware(), func(c *gin.Context) { middlewares.Role(c, "001") }, h.ProductReviewCreate)
	}
}
