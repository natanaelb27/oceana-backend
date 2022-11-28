package dtos

import (
	"strconv"
	"strings"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MerchantVoucherCreateRequest struct {
	Code                  string    `json:"code" binding:"required"`
	StartDatetime         time.Time `json:"start_datetime" binding:"required"`
	EndDatetime           time.Time `json:"end_datetime" binding:"required"`
	DiscountType          string    `json:"discount_type" binding:"required"`
	Discount              uint      `json:"discount" binding:"required"`
	MinimumSpendAmount    uint      `json:"minimum_spend_amount" binding:"required"`
	MaximumDiscountAmount uint      `json:"maximum_discount_amount"`
	MaximumQuota          uint      `json:"maximum_quota" binding:"required"`
}

func (source *MerchantVoucherCreateRequest) ToModel(id uint) (target *models.MerchantVoucher) {
	target = &models.MerchantVoucher{
		MerchantID:            id,
		Code:                  source.Code,
		StartDatetime:         source.StartDatetime,
		EndDatetime:           source.EndDatetime,
		DiscountType:          source.DiscountType,
		Discount:              source.Discount,
		MinimumSpendAmount:    source.MinimumSpendAmount,
		MaximumDiscountAmount: source.MaximumDiscountAmount,
		MaximumQuota:          source.MaximumQuota,
		CurrentQuota:          0,
	}
	return
}

type MerchantVoucherCreateResponse struct {
	ID                    uint      `json:"id"`
	Code                  string    `json:"code"`
	StartDatetime         time.Time `json:"start_datetime"`
	EndDatetime           time.Time `json:"end_datetime"`
	DiscountType          string    `json:"discount_type"`
	Discount              uint      `json:"discount"`
	MinimumSpendAmount    uint      `json:"minimum_spend_amount"`
	MaximumDiscountAmount uint      `json:"maximum_discount_amount"`
	MaximumQuota          uint      `json:"maximum_quota"`
	CurrentQuota          uint      `json:"current_quota"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func NewMerchantVoucherCreateResponse(source *models.MerchantVoucher) (target *MerchantVoucherCreateResponse) {
	target = &MerchantVoucherCreateResponse{
		ID:                    source.ID,
		Code:                  source.Code,
		StartDatetime:         source.StartDatetime,
		EndDatetime:           source.EndDatetime,
		DiscountType:          source.DiscountType,
		Discount:              source.Discount,
		MinimumSpendAmount:    source.MinimumSpendAmount,
		MaximumDiscountAmount: source.MaximumDiscountAmount,
		MaximumQuota:          source.MaximumQuota,
		CurrentQuota:          source.CurrentQuota,
		CreatedAt:             source.CreatedAt,
		UpdatedAt:             source.UpdatedAt,
	}
	return
}

type MerchantVoucherUpdateRequest struct {
	StartDatetime         time.Time `json:"start_datetime" binding:"required"`
	EndDatetime           time.Time `json:"end_datetime" binding:"required"`
	DiscountType          string    `json:"discount_type" binding:"required"`
	Discount              uint      `json:"discount" binding:"required"`
	MinimumSpendAmount    uint      `json:"minimum_spend_amount" binding:"required"`
	MaximumDiscountAmount uint      `json:"maximum_discount_amount"`
	MaximumQuota          uint      `json:"maximum_quota" binding:"required"`
}

func (source *MerchantVoucherUpdateRequest) ToModel(id uint) (target *models.MerchantVoucher) {
	target = &models.MerchantVoucher{
		ID:                    id,
		StartDatetime:         source.StartDatetime,
		EndDatetime:           source.EndDatetime,
		DiscountType:          source.DiscountType,
		Discount:              source.Discount,
		MinimumSpendAmount:    source.MinimumSpendAmount,
		MaximumDiscountAmount: source.MaximumDiscountAmount,
		MaximumQuota:          source.MaximumQuota,
	}
	return
}

type MerchantVoucherUpdateResponse struct {
	ID                    uint      `json:"id"`
	Code                  string    `json:"code"`
	StartDatetime         time.Time `json:"start_datetime"`
	EndDatetime           time.Time `json:"end_datetime"`
	DiscountType          string    `json:"discount_type"`
	Discount              uint      `json:"discount"`
	MinimumSpendAmount    uint      `json:"minimum_spend_amount"`
	MaximumDiscountAmount uint      `json:"maximum_discount_amount"`
	MaximumQuota          uint      `json:"maximum_quota"`
	CurrentQuota          uint      `json:"current_quota"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func NewMerchantVoucherUpdateResponse(source *models.MerchantVoucher) (target *MerchantVoucherUpdateResponse) {
	target = &MerchantVoucherUpdateResponse{
		ID:                    source.ID,
		Code:                  source.Code,
		StartDatetime:         source.StartDatetime,
		EndDatetime:           source.EndDatetime,
		DiscountType:          source.DiscountType,
		Discount:              source.Discount,
		MinimumSpendAmount:    source.MinimumSpendAmount,
		MaximumDiscountAmount: source.MaximumDiscountAmount,
		MaximumQuota:          source.MaximumQuota,
		CurrentQuota:          source.CurrentQuota,
		CreatedAt:             source.CreatedAt,
		UpdatedAt:             source.UpdatedAt,
	}
	return
}

type MerchantVoucherDeleteResponse struct {
	ID                    uint           `json:"id"`
	Code                  string         `json:"code"`
	StartDatetime         time.Time      `json:"start_datetime"`
	EndDatetime           time.Time      `json:"end_datetime"`
	DiscountType          string         `json:"discount_type"`
	Discount              uint           `json:"discount"`
	MinimumSpendAmount    uint           `json:"minimum_spend_amount"`
	MaximumDiscountAmount uint           `json:"maximum_discount_amount"`
	MaximumQuota          uint           `json:"maximum_quota"`
	CurrentQuota          uint           `json:"current_quota"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at"`
}

func NewMerchantVoucherDeleteResponse(source *models.MerchantVoucher) (target *MerchantVoucherDeleteResponse) {
	target = &MerchantVoucherDeleteResponse{
		ID:                    source.ID,
		Code:                  source.Code,
		StartDatetime:         source.StartDatetime,
		EndDatetime:           source.EndDatetime,
		DiscountType:          source.DiscountType,
		Discount:              source.Discount,
		MinimumSpendAmount:    source.MinimumSpendAmount,
		MaximumDiscountAmount: source.MaximumDiscountAmount,
		MaximumQuota:          source.MaximumQuota,
		CurrentQuota:          source.CurrentQuota,
		CreatedAt:             source.CreatedAt,
		UpdatedAt:             source.UpdatedAt,
		DeletedAt:             source.DeletedAt,
	}
	return
}

type MerchantVoucherGetByIDResponse struct {
	ID                    uint      `json:"id"`
	MerchantID            uint      `json:"merchant_id"`
	Code                  string    `json:"code"`
	StartDatetime         time.Time `json:"start_datetime"`
	EndDatetime           time.Time `json:"end_datetime"`
	DiscountType          string    `json:"discount_type"`
	Discount              uint      `json:"discount"`
	MinimumSpendAmount    uint      `json:"minimum_spend_amount"`
	MaximumDiscountAmount uint      `json:"maximum_discount_amount"`
	MaximumQuota          uint      `json:"maximum_quota"`
	CurrentQuota          uint      `json:"current_quota"`
	VoucherStatus         string    `json:"voucher_status"`
}

func NewMerchantVoucherGetByIDResponse(source *dtos.MerchantVoucherDetails) (target *MerchantVoucherGetByIDResponse) {
	target = &MerchantVoucherGetByIDResponse{
		ID:                    source.ID,
		MerchantID:            source.MerchantID,
		Code:                  source.Code,
		StartDatetime:         source.StartDatetime,
		EndDatetime:           source.EndDatetime,
		DiscountType:          source.DiscountType,
		Discount:              source.Discount,
		MinimumSpendAmount:    source.MinimumSpendAmount,
		MaximumDiscountAmount: source.MaximumDiscountAmount,
		MaximumQuota:          source.MaximumQuota,
		CurrentQuota:          source.CurrentQuota,
		VoucherStatus:         source.VoucherStatus,
	}
	return
}

type MerchantVoucherQuery struct {
	SearchVoucherCode   string `json:"search_voucher_code"`
	FilterVoucherStatus string `json:"filter_voucher_status"`
	MerchantID          uint   `json:"merchant_id"`
	SortBy              string `json:"sort_by"`
	Sort                string `json:"sort"`
	Limit               uint   `json:"limit"`
	Page                uint   `json:"page"`
}

func NewMerchantVoucherQuery(c *gin.Context) *MerchantVoucherQuery {
	var query MerchantVoucherQuery

	searchVoucherCode := c.Query("search_voucher_code")
	filterVoucherStatus := strings.ToLower(c.Query("filter_voucher_status"))
	merchantID := c.Query("merchant_id")
	sortBy := strings.ToLower(c.Query("sort_by"))
	sort := strings.ToLower(c.Query("sort"))
	limit := c.Query("limit")
	page := c.Query("page")

	query.SearchVoucherCode = searchVoucherCode
	query.FilterVoucherStatus = filterVoucherStatus

	merchantIDInt, err := strconv.Atoi(merchantID)

	if err != nil {
		query.MerchantID = 0
	} else if merchantIDInt <= 0 {
		query.MerchantID = 0
	} else {
		query.MerchantID = uint(merchantIDInt)
	}

	if sortBy == "end_datetime" {
		query.SortBy = "end_datetime"
	} else {
		query.SortBy = "start_datetime"
	}

	if sort == "asc" {
		query.Sort = "ASC"
	} else {
		query.Sort = "DESC"
	}

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

type MerchantVoucherGetAllResponse struct {
	Pagination      *utils.Pagination              `json:"pagination"`
	MerchantVoucher []*dtos.MerchantVoucherDetails `json:"merchant_voucher"`
}

func NewMerchantVoucherGetAllResponse(source []*dtos.MerchantVoucherDetails, pagination *utils.Pagination) (target *MerchantVoucherGetAllResponse) {
	if pagination == nil {
		return nil
	}

	target = &MerchantVoucherGetAllResponse{
		Pagination:      pagination,
		MerchantVoucher: source,
	}

	return
}
