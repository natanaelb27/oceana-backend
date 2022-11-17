package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) MerchantVoucherCreate(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.MerchantVoucherCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	merchantVouherCreateResponse, err := h.merchantVoucherService.Create(id, request)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, merchantVouherCreateResponse, nil)
}
