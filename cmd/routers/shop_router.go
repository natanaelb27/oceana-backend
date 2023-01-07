package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/middlewares"
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
	privateRoute := private.Group(path)
	{
		privateRoute.GET("/vouchers/buyer", func(c *gin.Context) { middlewares.Role(c, "001") }, h.MerchantVoucherGetAllByBuyer)
	}
}
