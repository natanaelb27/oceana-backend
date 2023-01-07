package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
)

func (h *Handler) MarketplaceVoucherGetAll(c *gin.Context) {
	query := dtos.NewMarketplaceVoucherQuery(c)
	results, pagination, err := h.marketplaceVoucherService.GetAll(query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	var responseData []*dtos.MarketplaceVoucherGetResponse = make([]*dtos.MarketplaceVoucherGetResponse, 0)

	for _, result := range results {
		responseData = append(responseData, dtos.NewMarketplaceVoucherGetResponse(result))
	}

	utils.ResponseSuccess(c, responseData, pagination)
}

func (h *Handler) MarketplaceVoucherGetAllByActive(c *gin.Context) {
	results, err := h.marketplaceVoucherService.GetAllByActive()
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	var responseData []*dtos.MarketplaceVoucherGetResponse = make([]*dtos.MarketplaceVoucherGetResponse, 0)

	for _, result := range results {
		responseData = append(responseData, dtos.NewMarketplaceVoucherGetResponse(result))
	}

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) MarketplaceVoucherGetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.marketplaceVoucherService.GetById(uint(id))
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewMarketplaceVoucherGetResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) MarketplaceVoucherCreate(c *gin.Context) {
	var request *dtos.MarketplaceVoucherCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	model := request.ToModel()

	result, err := h.marketplaceVoucherService.Create(model)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewMarketplaceVoucherCreateResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) MarketplaceVoucherUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}
	var request *dtos.MarketplaceVoucherCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	model := request.ToModel()

	result, err := h.marketplaceVoucherService.Update(uint(id), model)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewMarketplaceVoucherCreateResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) MarketplaceVoucherDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.marketplaceVoucherService.Delete(uint(id))
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	responseData := dtos.NewMarketplaceVoucherDeleteResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) MarketplaceVoucherGetAllByBuyer(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	query := dtos.NewBuyerMarketplaceVoucherQuery(c)
	merchantVouherGetAllResponse, pagination, err := h.marketplaceVoucherService.GetAllByBuyer(id, query)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, merchantVouherGetAllResponse, pagination)
}
