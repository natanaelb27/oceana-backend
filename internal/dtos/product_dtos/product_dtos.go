package product_dtos

import (
	"strconv"
	"strings"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Merchant struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

type Product struct {
	ID         uint `json:"id"`
	MerchantID uint `json:"merchant_id"`
	// Merchant      models.Merchant `json:"merchant"`
	CategoryID uint `json:"category_id"`
	// Category      models.Category `json:"category"`
	Name                      string    `json:"name"`
	Description               string    `json:"description"`
	MediaURL                  string    `json:"media_url"`
	IsHazardous               bool      `json:"is_hazardous"`
	Weight                    uint      `json:"weight"`
	PackageSize               string    `json:"package_size"`
	Condition                 string    `json:"condition"`
	ProductStatus             string    `json:"product_status"`
	PriceMinimum              uint      `json:"price_minimum"`
	PriceMaximum              uint      `json:"price_maximum"`
	PriceMinimumAfterDiscount float64   `json:"price_minimum_after_discount"`
	PriceMaximumAfterDiscount float64   `json:"price_maximum_after_discount"`
	Discount                  uint      `json:"discount"`
	TotalViewCount            uint      `json:"total_view_count"`
	Merchant                  Merchant  `json:"merchant"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

type ProductQuery struct {
	Search string
	SortBy string
	Sort   string
	Limit  uint
	Page   uint
	// Filter
	MerchantID      uint
	CategoryID      uint
	PriceMinimum    uint
	PriceMaximum    uint
	MerchantCities  []string
	CourierServices []string
}

func NewProductQuery(c *gin.Context) *ProductQuery {
	query := &ProductQuery{}

	query.Search = c.Query("search")
	query.MerchantID = utils.StringToUint(c.Query("merchant-id"))
	query.CategoryID = utils.StringToUint(c.Query("category-id"))
	query.PriceMinimum = utils.StringToUint(c.Query("price-minimum"))
	query.PriceMaximum = utils.StringToUint(c.Query("price-maximum"))

	query.Page = utils.StringToUint(c.Query("page"))
	if query.Page == 0 {
		query.Page = 1
	}

	query.Limit = utils.StringToUint(c.Query("limit"))
	if query.Limit == 0 {
		query.Limit = 10
	}

	merchantCitiesString := c.Query("merchant-cities")
	if merchantCitiesString != "" {
		query.MerchantCities = strings.Split(merchantCitiesString, ",")
	}

	courierServicesString := c.Query("courier-services")
	if courierServicesString != "" {
		query.CourierServices = strings.Split(courierServicesString, ",")
	}

	query.Sort = utils.QueryParamSort(c.Query("sort"))

	sortBy := strings.ToLower(c.Query("sort-by"))
	switch sortBy {
	case "newest":
		query.SortBy = "products.created_at"
	case "price":
		query.SortBy = "price_minimum"
	case "recommended":
		query.SortBy = "total_view_count"
	// case most buy
	default:
		query.SortBy = "products.name"
	}

	return query
}

func NewShopProductQuery(c *gin.Context) (query *ProductQuery, err error) {
	query = NewProductQuery(c)
	merchantId, err := strconv.Atoi(c.Param("merchantId"))
	if err != nil {
		err = e.ErrInternalOperationFailed
		return
	}

	query.MerchantID = uint(merchantId)

	return
}
