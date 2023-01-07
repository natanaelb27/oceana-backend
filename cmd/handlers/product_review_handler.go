package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

func (h *Handler) ProductReviewCreate(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}
	orderDetailID, err := uuid.Parse(c.Param("orderDetailID"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidOrderDetailID)
		return
	}
	var request *dtos.ProductReviewCreateRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}
	request.UserID = id
	request.OrderDetailID = orderDetailID

	tx := c.MustGet("tx").(*gorm.DB)

	response, err := h.productReviewService.WithTx(tx).Create(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, response, nil)
}

func (h *Handler) ProductReviewGet(c *gin.Context) {
	productID := utils.StringToUint(c.Param("productID"))
	query := dtos.NewProductReviewQuery(c)
	response, pagination, err := h.productReviewService.GetByProductID(productID, query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, response, pagination)
}
