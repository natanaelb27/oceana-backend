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
		publicRoute.GET("", h.ProductGetAll)
		publicRoute.GET("/:product-id", h.ProductDetailsGetByID)
	}
	privateRoute := private.Group(path)
	{
		privateRoute.GET("/purchased-quantities", func(c *gin.Context) { middlewares.Role(c, "001") }, h.GetPurchaseQuantitiesByID)
		privateRoute.POST("/:product-id/view", func(c *gin.Context) { middlewares.Role(c, "011") }, h.ViewProduct)
		privateRoute.GET("/:product-id/favorite", func(c *gin.Context) { middlewares.Role(c, "011") }, h.GetFavoriteProduct)
		privateRoute.POST("/:product-id/favorite", func(c *gin.Context) { middlewares.Role(c, "011") }, h.FavoriteProduct)
		publicRoute.GET("/:product-id/model", h.ProductModelById)
		privateRoute.GET("/available-products", h.GetAvailableProductsByPromotionPeriod)
	}
}
