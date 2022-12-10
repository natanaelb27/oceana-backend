package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CitiesGetAll(c *gin.Context) {
	results, err := h.cityService.GetAll()
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, nil)
}
