package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func merchant(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/merchant"
	publicRoute := public.Group(path)
	{
		publicRoute.GET("/profile/:merchantId", h.MerchantGetInformationByID)
		publicRoute.GET("/vouchers", h.MerchantVoucherGetAll)
		publicRoute.GET("/vouchers/:merchantVoucherId", h.MerchantVoucherGetByID)
	}
	privateRoute := private.Group(path)
	{
		privateRoute.POST("/vouchers", h.MerchantVoucherCreate)
		privateRoute.PUT("/vouchers/:merchantVoucherId", h.MerchantVoucherUpdate)
		privateRoute.DELETE("/vouchers/:merchantVoucherId", h.MerchantVoucherDelete)
	}
}
