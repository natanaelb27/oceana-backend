package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) AddressCreate(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.AddressCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	model := request.ToModel(id)

	result, err := h.addressService.Create(model)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewAddressCreateResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}
