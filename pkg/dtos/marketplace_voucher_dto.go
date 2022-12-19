package dtos

import (
	"strconv"
	"strings"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type MarketplaceVoucherCreateRequest struct {
	Name                  string    `json:"name" binding:"required"`
	StartDatetime         time.Time `json:"start_datetime" binding:"required"`
	EndDatetime           time.Time `json:"end_datetime" binding:"required"`
	DiscountType          string    `json:"discount_type" binding:"required"`
	Discount              uint      `json:"discount" binding:"required"`
	MinimumSpendAmount    uint      `json:"minimum_spend_amount"`
	MaximumDiscountAmount uint      `json:"maximum_discount_amount"`
	MaximumQuota          uint      `json:"maximum_quota" binding:"required"`
}

func (source *MarketplaceVoucherCreateRequest) ToModel() (target *models.MarketplaceVoucher) {
	target = &models.MarketplaceVoucher{
		Name:                  source.Name,
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

type MarketplaceVoucherCreateResponse struct {
	ID                    uint      `json:"id"`
	Name                  string    `json:"name"`
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

func NewMarketplaceVoucherCreateResponse(source *models.MarketplaceVoucher) (target *MarketplaceVoucherCreateResponse) {
	target = &MarketplaceVoucherCreateResponse{
		ID:                    source.ID,
		Name:                  source.Name,
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

type MarketplaceVoucherGetResponse struct {
	ID                    uint      `json:"id"`
	Name                  string    `json:"name"`
	StartDatetime         time.Time `json:"start_datetime"`
	EndDatetime           time.Time `json:"end_datetime"`
	DiscountType          string    `json:"discount_type"`
	Discount              uint      `json:"discount"`
	MinimumSpendAmount    uint      `json:"minimum_spend_amount"`
	MaximumDiscountAmount uint      `json:"maximum_discount_amount"`
	MaximumQuota          uint      `json:"maximum_quota"`
	CurrentQuota          uint      `json:"current_quota"`
	VoucherStatus         string    `json:"voucher_status"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func NewMarketplaceVoucherGetResponse(source *models.MarketplaceVoucher) (target *MarketplaceVoucherGetResponse) {
	target = &MarketplaceVoucherGetResponse{
		ID:                    source.ID,
		Name:                  source.Name,
		StartDatetime:         source.StartDatetime,
		EndDatetime:           source.EndDatetime,
		DiscountType:          source.DiscountType,
		Discount:              source.Discount,
		MinimumSpendAmount:    source.MinimumSpendAmount,
		MaximumDiscountAmount: source.MaximumDiscountAmount,
		MaximumQuota:          source.MaximumQuota,
		CurrentQuota:          source.CurrentQuota,
		VoucherStatus:         utils.PromotionOrVoucherStatus(source.StartDatetime, source.EndDatetime),
		CreatedAt:             source.CreatedAt,
		UpdatedAt:             source.UpdatedAt,
	}

	return
}

type MarketplaceVoucherDeleteResponse struct {
	ID uint `json:"id"`
}

func NewMarketplaceVoucherDeleteResponse(source *models.MarketplaceVoucher) (target *MarketplaceVoucherDeleteResponse) {
	target = &MarketplaceVoucherDeleteResponse{
		ID: source.ID,
	}

	return
}

type MarketplaceVoucherGetAllResponse struct {
	MarketplaceVoucher []*dtos.MarketplaceVoucherDetails `json:"merchant_voucher"`
}

func NewMarketplaceVoucherGetAllResponse(source []*dtos.MarketplaceVoucherDetails) (target *MarketplaceVoucherGetAllResponse) {
	target = &MarketplaceVoucherGetAllResponse{
		MarketplaceVoucher: source,
	}

	return
}

type BuyerMarketplaceVoucherQuery struct {
	FilterVoucherStatus       string `json:"filter-voucher-status"`
	FilterVoucherAvailability string `json:"filter-voucher-availability"`
	SortBy                    string `json:"sort-by"`
	Sort                      string `json:"sort"`
	Limit                     uint   `json:"limit"`
	Page                      uint   `json:"page"`
}

func NewBuyerMarketplaceVoucherQuery(c *gin.Context) *BuyerMarketplaceVoucherQuery {
	var query BuyerMarketplaceVoucherQuery

	filterVoucherStatus := strings.ToLower(c.Query("filter-voucher-status"))
	filterVoucherAvailability := strings.ToLower(c.Query("filter-voucher-availability"))
	sortBy := strings.ToLower(c.Query("sort-by"))
	sort := strings.ToLower(c.Query("sort"))
	limit := c.Query("limit")
	page := c.Query("page")

	query.FilterVoucherStatus = filterVoucherStatus
	query.FilterVoucherAvailability = filterVoucherAvailability

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

type MarketplaceVoucherQuery struct {
	SearchVoucherName   string `json:"search-voucher-name"`
	FilterVoucherStatus string `json:"filter-voucher-status"`
	SortBy              string `json:"sort-by"`
	Sort                string `json:"sort"`
	Limit               uint   `json:"limit"`
	Page                uint   `json:"page"`
}

func NewMarketplaceVoucherQuery(c *gin.Context) *MarketplaceVoucherQuery {
	var query MarketplaceVoucherQuery

	searchVoucherName := c.Query("search-voucher-name")
	filterVoucherStatus := strings.ToLower(c.Query("filter-voucher-status"))
	sortBy := strings.ToLower(c.Query("sort-by"))
	sort := strings.ToLower(c.Query("sort"))
	limit := c.Query("limit")
	page := c.Query("page")

	query.SearchVoucherName = searchVoucherName
	query.FilterVoucherStatus = filterVoucherStatus

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
