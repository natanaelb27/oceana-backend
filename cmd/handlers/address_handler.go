package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AddressCreate(c *gin.Context) {
	var request *dtos.AddressCreateRequest
	userContext := c.MustGet("user").(models.User)

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	model := request.ToModel(userContext.ID)

	result, err := h.addressService.Create(model)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewAddressCreateResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}
