package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func review(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/reviews"
	privateRoute := private.Group(path)
	{
		privateRoute.POST("/order-details/:orderDetailID", h.ProductReviewCreate)
	}
}
