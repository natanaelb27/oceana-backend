package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
)

func (h *Handler) GetRajaOngkirShippingCost(c *gin.Context) {
	var request *dtos.GetRajaOngkirShippingCostRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	result, err := h.rajaOngkirService.GetRajaOngkirShippingCost(request)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, result, nil)
}
