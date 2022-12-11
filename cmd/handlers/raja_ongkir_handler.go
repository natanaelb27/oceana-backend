package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/gin-gonic/gin"
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
