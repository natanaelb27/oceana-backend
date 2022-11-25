package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func review(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/reviews"
	publicRoute := public.Group(path) 
	{
		publicRoute.GET("products/:productID", h.ProductReviewGet)
	}
	privateRoute := private.Group(path)
	{
		privateRoute.POST("/order-details/:orderDetailID", h.ProductReviewCreate)
	}
}
