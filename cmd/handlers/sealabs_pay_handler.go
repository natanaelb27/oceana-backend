package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) SealabsPayCreate(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.SealabsPayCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	sealabsPayCreateResponse, err := h.sealabsPayService.Create(id, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, sealabsPayCreateResponse, nil)
}

func (h *Handler) SealabsPayGetAll(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	results, err := h.sealabsPayService.GetAll(id)
	if err != nil {
		utils.ResponseError(c, err)
	}

	var responseData []*dtos.SealabsPayCreateResponse = make([]*dtos.SealabsPayCreateResponse, 0)

	for _, result := range results {
		responseData = append(responseData, dtos.NewSealabsPayCreateResponse(result))
	}

	utils.ResponseSuccess(c, responseData, nil)

}

func (h *Handler) SealabsPayGetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("sealabs_pay_id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.sealabsPayService.GetById(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewSealabsPayCreateResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) SealabsPayDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("sealabs_pay_id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.sealabsPayService.Delete(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewSealabsPayDeleteResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}
