package handlers

import (
	"fmt"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) MerchantPromotionCreate(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var mpRequest *dtos.MerchantPromotionCreateRequest
	if err := c.ShouldBindJSON(&mpRequest); err != nil {
		utils.ResponseError(c, err)
		return
	}

	fmt.Println(mpRequest)
	responses, err := h.merchantPromotionService.CreateMyShopPromotion(id, mpRequest, mpRequest.ProductMerchantPromotion)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, responses, nil)
}
