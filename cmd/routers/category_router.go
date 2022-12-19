package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
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
