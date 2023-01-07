package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
)

func category(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/categories"
	publicRoute := public.Group(path)
	{
		publicRoute.GET("", h.CategoryGetAll)
		publicRoute.GET("/ancestors/:category-id", h.CategoryAncestorsGetByID)
		publicRoute.GET("/merchants/:merchant-id", h.GetMerchantCategories)
	}
}
