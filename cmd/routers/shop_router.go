package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func shop(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/shop"
	publicRoute := public.Group(path)
	{
		publicRoute.GET("/information/:merchant-id", h.MerchantGetInformationByID)
		publicRoute.GET("/vouchers", h.MerchantVoucherGetAll)
		publicRoute.GET("/vouchers/:merchant-voucher-id", h.MerchantVoucherGetByID)
		publicRoute.GET("/:merchantId/products", h.ShopProductGetAll)

	}
}
