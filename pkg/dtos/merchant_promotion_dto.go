package dtos

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/internal/models"
	"github.com/natanaelb27/oceana-backend/internal/utils"
)

type MerchantPromotionRequest struct {
	Name                     string                             `json:"name" binding:"required"`
	StartDatetime            time.Time                          `json:"start_datetime" binding:"required"`
	EndDatetime              time.Time                          `json:"end_datetime" binding:"required"`
	ProductIDList            []int64                            `json:"product_id_list" binding:"required"`
	ProductMerchantPromotion []*ProductMerchantPromotionRequest `json:"product_merchant_promotion" binding:"required"`
}

func (source *MerchantPromotionRequest) ToModel(merchantID uint) (target *models.MerchantPromotion) {
	target = &models.MerchantPromotion{
		MerchantId:    merchantID,
		Name:          source.Name,
		StartDatetime: source.StartDatetime,
		EndDatetime:   source.EndDatetime,
	}

	return
}

type MerchantPromotionResponse struct {
	ID                       uint                                `json:"id"`
	MerchantId               uint                                `json:"merchant_id"`
	Name                     string                              `json:"name"`
	StartDatetime            time.Time                           `json:"start_datetime"`
	EndDatetime              time.Time                           `json:"end_datetime"`
	PromotionStatus          string                              `json:"promotion_status"`
	CreatedAt                time.Time                           `json:"created_at"`
	UpdatedAt                time.Time                           `json:"updated_at"`
	ProductMerchantPromotion []*ProductMerchantPromotionResponse `json:"product_merchant_promotion"`
}

func NewMerchantPromotionResponse(source *models.MerchantPromotion, pmpResponse []*ProductMerchantPromotionResponse) (target *MerchantPromotionResponse) {
	target = &MerchantPromotionResponse{
		ID:                       source.ID,
		MerchantId:               source.MerchantId,
		Name:                     source.Name,
		StartDatetime:            source.StartDatetime,
		EndDatetime:              source.EndDatetime,
		PromotionStatus:          utils.PromotionOrVoucherStatus(source.StartDatetime, source.EndDatetime),
		CreatedAt:                source.CreatedAt,
		UpdatedAt:                source.UpdatedAt,
		ProductMerchantPromotion: pmpResponse,
	}

	return
}

type MerchantPromotionQuery struct {
	SearchPromotionName   string `json:"search_promotion_name"`
	FilterPromotionStatus string `json:"filter_promotion_status"`
	Limit                 uint   `json:"limit"`
	Page                  uint   `json:"page"`
}

func NewMerchantPromotionQuery(c *gin.Context) *MerchantPromotionQuery {
	var query MerchantPromotionQuery

	searchPromotionName := c.Query("search-promotion-name")
	filterPromotionStatus := c.Query("filter-promotion-status")
	limit := c.Query("limit")
	page := c.Query("page")

	query.SearchPromotionName = searchPromotionName
	query.FilterPromotionStatus = filterPromotionStatus

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		query.Page = 1
	} else if pageInt <= 0 {
		query.Page = 1
	} else {
		query.Page = uint(pageInt)
	}

	limitInt, err := strconv.Atoi(limit)

	if err != nil {
		query.Limit = 10
	} else if limitInt <= 0 {
		query.Limit = 10
	} else {
		query.Limit = uint(limitInt)
	}

	return &query
}

type MerchantPromotionDetails struct {
	ID                       uint                                `json:"id"`
	MerchantId               uint                                `json:"merchant_id"`
	Name                     string                              `json:"name"`
	StartDatetime            time.Time                           `json:"start_datetime"`
	EndDatetime              time.Time                           `json:"end_datetime"`
	PromotionStatus          string                              `json:"promotion_status"`
	ProductMerchantPromotion []*ProductMerchantPromotionResponse `json:"product_merchant_promotion"`
}

func NewMerchantPromotionDetails(source *models.MerchantPromotion) (target *MerchantPromotionDetails) {
	target = &MerchantPromotionDetails{
		ID:              source.ID,
		MerchantId:      source.MerchantId,
		Name:            source.Name,
		StartDatetime:   source.StartDatetime,
		EndDatetime:     source.EndDatetime,
		PromotionStatus: utils.PromotionOrVoucherStatus(source.StartDatetime, source.EndDatetime),
	}

	return
}
