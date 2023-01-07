package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
)

func (h *Handler) SubDistrictGetAll(c *gin.Context) {
	query := dtos.NewSubDistrictQuery(c)

	results, err := h.subDistrictService.GetAll(query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, nil)
}
