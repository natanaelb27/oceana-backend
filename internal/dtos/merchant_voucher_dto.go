package dtos

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type MerchantVoucherDetails struct {
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
