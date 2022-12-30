package handlers

import (
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos/product_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (h *Handler) ProductDetailsGetByID(c *gin.Context) {
	productIDUint64, err := strconv.ParseUint(c.Param("product-id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	internalQuery := &internalDtos.ProductDetailsInternalQuery{
		ProductID: uint(productIDUint64),
		Status:    constants.ProductStatusActive,
	}

	productDetailsGetByIDResponse, err := h.productService.ProductDetailsGetByID(internalQuery)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, productDetailsGetByIDResponse, nil)

}

func (h *Handler) ProductGetAll(c *gin.Context) {
	query := product_dtos.NewProductQuery(c)

	results, pagination, err := h.productService.GetAll(query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, pagination)
}

func (h *Handler) MyShopGetAll(c *gin.Context) {
	query := product_dtos.NewProductQuery(c)

	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}
	merchant, err := h.merchantService.MerchantGetProfile(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	query.MerchantID = merchant.ID

	results, pagination, err := h.productService.MyShopGetAll(query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, pagination)
}

func (h *Handler) GetPurchaseQuantitiesByID(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	query := dtos.NewPurchaseQuantitiesQuery(c)

	getPurchaseQuantitiesByIDResponse, err := h.productService.GetPurchaseQuantitiesByID(id, query)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, getPurchaseQuantitiesByIDResponse, nil)
}

func (h *Handler) ShopProductGetAll(c *gin.Context) {
	query, err := product_dtos.NewShopProductQuery(c)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	results, pagination, err := h.productService.GetAll(query)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, pagination)
}

func (h *Handler) ViewProduct(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	productIDUint64, err := strconv.ParseUint(c.Param("product-id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	productID := uint(productIDUint64)

	viewProductResponse, err := h.productService.ViewProduct(id, productID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, viewProductResponse, nil)
}

func (h *Handler) GetFavoriteProduct(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	productIDUint64, err := strconv.ParseUint(c.Param("product-id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	productID := uint(productIDUint64)

	getFavoriteProductResponse, err := h.productService.GetFavoriteProduct(id, productID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, getFavoriteProductResponse, nil)
}

func (h *Handler) FavoriteProduct(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	productIDUint64, err := strconv.ParseUint(c.Param("product-id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	productID := uint(productIDUint64)

	getFavoriteProductResponse, err := h.productService.FavoriteProduct(id, productID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, getFavoriteProductResponse, nil)
}

func (h *Handler) ProductCreate(c *gin.Context) {
	var request *product_dtos.ProductModel
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	trx := c.MustGet("tx").(*gorm.DB)
	response, err := h.productService.WithTx(trx).CreateProduct(request, id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, response, nil)
}

func (h *Handler) ProductCheckCourierService(c *gin.Context) {
	var request *product_dtos.ProductCheckCourierService
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	response, err := h.productService.CheckCourierService(request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, response, nil)
}

func (h *Handler) ProductUpdate(c *gin.Context) {
	productIDInt, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidParams)
		return
	}

	userId := utils.GinGetUUID(c, "id")
	if userId == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *product_dtos.ProductModel
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	request.ID = uint(productIDInt)

	trx := c.MustGet("tx").(*gorm.DB)
	response, err := h.productService.WithTx(trx).UpdateProductUsingModel(userId, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, response, nil)
}

func (h *Handler) ProductModelById(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("product-id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	productModel, err := h.productService.GetUpdateModelByProductID(uint(productID))
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, productModel, nil)
}

func (h *Handler) GetAvailableProductsByPromotionPeriod(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	query := dtos.NewAvailableProductsByPromotionPeriodQuery(c)

	result, err := h.productService.GetAvailableProductsByPromotionPeriod(id, query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, result, nil)
}

func (h *Handler) ProductLogisticsByMerchantID(c *gin.Context) {
	merchantID, err := strconv.Atoi(c.Param("merchant-id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidParams)
		return
	}

	results, err := h.logisticService.ProductGetByMerchantID(uint(merchantID))
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, nil)
}

func (h *Handler) GetAllUserFavoriteProducts(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	query := product_dtos.NewProductFavoriteQuery(c)

	results, pagination, err := h.productService.GetAllUserFavoriteProducts(id, query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, pagination)
}

func (h *Handler) DeleteFavoriteProducts(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.DeleteFavoriteProductsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	tx := c.MustGet("tx").(*gorm.DB)

	results, err := h.productService.WithTx(tx).DeleteFavoriteProducts(id, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, nil)
}
