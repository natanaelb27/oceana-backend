package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
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
		privateRoute.POST("/check-courier-service", func(c *gin.Context) { middlewares.Role(c, "001") }, h.ProductCheckCourierService)
		privateRoute.GET("/logistics/:merchant-id", func(c *gin.Context) { middlewares.Role(c, "001") }, h.ProductLogisticsByMerchantID)
		privateRoute.GET("/purchased-quantities", func(c *gin.Context) { middlewares.Role(c, "001") }, h.GetPurchaseQuantitiesByID)
		privateRoute.POST("/:product-id/view", func(c *gin.Context) { middlewares.Role(c, "011") }, h.ViewProduct)
		privateRoute.GET("/:product-id/favorite", func(c *gin.Context) { middlewares.Role(c, "011") }, h.GetFavoriteProduct)
		privateRoute.POST("/:product-id/favorite", func(c *gin.Context) { middlewares.Role(c, "011") }, h.FavoriteProduct)
		privateRoute.GET("/:product-id/model", func(c *gin.Context) { middlewares.Role(c, "010") }, h.ProductModelById)
		privateRoute.GET("/available-products", func(c *gin.Context) { middlewares.Role(c, "010") }, h.GetAvailableProductsByPromotionPeriod)
	}
}
