package handlers

import (
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ProductCreate(c *gin.Context) {
	merchantIDUint64, err := strconv.ParseUint(c.Param("merchant_id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidParams)
		return
	}

	merchantID := uint(merchantIDUint64)

	var request *dtos.ProductCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	model := request.ToModel(merchantID)

	result, err := h.productService.CreateProduct(model)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewProductCreateResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}
