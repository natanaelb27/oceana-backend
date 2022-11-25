package handlers

import (
	in_cart_dtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/cart_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CartGetAllByUserId(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")

	results, err := h.cartService.GetAllByUserId(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	var responseData []*in_cart_dtos.Merchant = make([]*in_cart_dtos.Merchant, 0)

	responseData = append(responseData, results...)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) CartUpdate(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")

	var request *dtos.CartUpsertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	internalRequest := request.CartUpsertRequestToInternalCartRequest()

	internalRequest.UserID = id

	result, err := h.cartService.Update(internalRequest)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewCartResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}
