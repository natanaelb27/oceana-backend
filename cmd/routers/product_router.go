package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-gonic/gin"
)

func product(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/products"
	publicRoute := public.Group(path)
	{
		publicRoute.GET("/:productId", h.ProductDetailsGetByID)
	}
	privateRoute := private.Group(path)
	{
		privateRoute.GET("/purchased-quantities", func(c *gin.Context) { middlewares.Role(c, "001") }, h.GetPurchaseQuantitiesByID)
	}
}
