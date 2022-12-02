package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CourierGetAll(c *gin.Context) {
	results, err := h.courierService.GetAll()
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, nil)
}
