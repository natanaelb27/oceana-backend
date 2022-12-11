package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CitiesGetAll(c *gin.Context) {
	query := dtos.NewCitiesQuery(c)

	results, err := h.cityService.GetAll(query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, nil)
}
