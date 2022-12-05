package handlers

import (
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) MerchantGetInformationByID(c *gin.Context) {
	merchantIDUint64, err := strconv.ParseUint(c.Param("merchant-id"), 10, 64)
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

func (h *Handler) MerchantGetProfile(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}
	responseData, err := h.merchantService.MerchantGetProfile(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, responseData, nil)
}
