package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type MerchantPromotionCreateRequest struct {
	Name                     string                                   `json:"name" binding:"required"`
	StartDatetime            time.Time                                `json:"start_date_time" binding:"required"`
	EndDatetime              time.Time                                `json:"end_date_time" binding:"required"`
	ProductIDList            []int64                                  `json:"product_id_list" binding:"required"`
	ProductMerchantPromotion []*ProductMerchantPromotionCreateRequest `json:"product_merchant_promotion" binding:"required"`
}

func (source *MerchantPromotionCreateRequest) ToModel(merchantID uint) (target *models.MerchantPromotion) {
	target = &models.MerchantPromotion{
		MerchantId:    merchantID,
		Name:          source.Name,
		StartDatetime: source.StartDatetime,
		EndDatetime:   source.EndDatetime,
	}

	return
}

type MerchantPromotionCreateResponse struct {
	ID            uint      `json:"id"`
	MerchantId    uint      `json:"merchant_id"`
	Name          string    `json:"name"`
	StartDatetime time.Time `json:"start_date_time"`
	EndDatetime   time.Time `json:"end_date_time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewMerchantPromotionCreateResponse(source *models.MerchantPromotion) (target *MerchantPromotionCreateResponse) {
	target = &MerchantPromotionCreateResponse{
		ID:            source.ID,
		MerchantId:    source.MerchantId,
		Name:          source.Name,
		StartDatetime: source.StartDatetime,
		EndDatetime:   source.EndDatetime,
		CreatedAt:     source.CreatedAt,
		UpdatedAt:     source.UpdatedAt,
	}

	return
}

type MerchantPromotionQuery struct {
	SearchPromotionName   string `json:"search_promotion_name"`
	FilterPromotionStatus string `json:"filter_promotion_status"`
}

func NewMerchantPromotionQuery(c *gin.Context) *MerchantPromotionQuery {
	var query MerchantPromotionQuery

	searchPromotionName := c.Query("search-promotion-name")
	filterPromotionStatus := c.Query("filter-promotion-status")

	query.SearchPromotionName = searchPromotionName
	query.FilterPromotionStatus = filterPromotionStatus

	return &query
}

type MerchantPromotionDetails struct {
	ID              uint      `json:"id"`
	MerchantId      uint      `json:"merchant_id"`
	Name            string    `json:"name"`
	StartDatetime   time.Time `json:"start_date_time"`
	EndDatetime     time.Time `json:"end_date_time"`
	PromotionStatus string    `json:"promotion_status"`
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

type MerchantPromotionUpdateRequest struct {
	Name                     string                                   `json:"name" binding:"required"`
	StartDatetime            time.Time                                `json:"start_date_time" binding:"required"`
	EndDatetime              time.Time                                `json:"end_date_time" binding:"required"`
	ProductIDList            []int64                                  `json:"product_id_list" binding:"required"`
	ProductMerchantPromotion []*ProductMerchantPromotionCreateRequest `json:"product_merchant_promotion" binding:"required"`
}

type MerchantPromotionUpdateResponse struct {
	ID            uint      `json:"id"`
	MerchantId    uint      `json:"merchant_id"`
	Name          string    `json:"name"`
	StartDatetime time.Time `json:"start_date_time"`
	EndDatetime   time.Time `json:"end_date_time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
