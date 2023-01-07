package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
)

func rajaOngkir(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/raja-ongkir"
	publicRoute := public.Group(path)
	{
		publicRoute.POST("/shipping-cost", h.GetRajaOngkirShippingCost)
	}
}
