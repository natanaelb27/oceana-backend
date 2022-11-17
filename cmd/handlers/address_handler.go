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

func (h *Handler) AddressGetAll(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
	}

	results, err := h.addressService.GetProfileAddressess(id)
	if err != nil {
		utils.ResponseError(c, err)
	}

	var responseData []*dtos.AddressGetResponse = make([]*dtos.AddressGetResponse, 0)

	for _, result := range results {
		responseData = append(responseData, dtos.NewAddressGetResponse(result))
	}

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) AddressGetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("address_id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.addressService.GetProfileAddressById(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewAddressGetResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) AddressUpdate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("address_id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	userID := utils.GinGetUUID(c, "id")
	if userID == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
	}

	var request *dtos.AddressCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	model := request.ToModel(userID)

	result, err := h.addressService.UpdateProfileAddress(id, model)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewAddressCreateResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}
