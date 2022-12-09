package routers

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"github.com/gin-gonic/gin"
)

func rajaOngkir(public *gin.RouterGroup, private *gin.RouterGroup, h *handlers.Handler) {
	path := "/raja-ongkir"
	publicRoute := public.Group(path)
	{
		publicRoute.POST("/shipping-cost", h.GetRajaOngkirShippingCost)
	}
}
