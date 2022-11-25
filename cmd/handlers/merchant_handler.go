package handlers

import (
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) MerchantGetInformationByID(c *gin.Context) {
	merchantIDUint64, err := strconv.ParseUint(c.Param("merchantId"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	merchantID := uint(merchantIDUint64)

	merchantGetInformationByIDResponse, err := h.merchantService.GetInformationByID(merchantID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, merchantGetInformationByIDResponse, nil)
}
