package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/internal/utils"
)

func (h *Handler) CourierGetAll(c *gin.Context) {
	results, err := h.courierService.GetAll()
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, nil)
}
