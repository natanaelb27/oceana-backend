package product_dtos

import (
	"fmt"
	"strings"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type Product struct {
	ID         uint `json:"id"`
	MerchantID uint `json:"merchant_id"`
	// Merchant      models.Merchant `json:"merchant"`
	CategoryID uint `json:"category_id"`
	// Category      models.Category `json:"category"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	IsHazardous    bool      `json:"is_hazardous"`
	Weight         uint      `json:"weight"`
	PackageSize    string    `json:"package_size"`
	Condition      string    `json:"condition"`
	ProductStatus  string    `json:"product_status"`
	PriceMinimum   uint      `json:"price_minimum"`
	PriceMaximum   uint      `json:"price_maximum"`
	TotalViewCount uint      `json:"total_view_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ProductQuery struct {
	Search string `json:"search"`
	SortBy string `json:"sort_by"`
	Sort   string `json:"sort"`
	Limit  uint   `json:"limit"`
	Page   uint   `json:"page"`
	// Filter
	MerchantID   uint `json:"merchant_id"`
	PriceMinimum uint `json:"price_minimum"`
	PriceMaximum uint `json:"price_maximum"`
}

func NewProductQuery(c *gin.Context) *ProductQuery {
	query := &ProductQuery{}

	query.Search = c.Query("search")
	query.MerchantID = utils.StringToUint(c.Query("merchant_id"))
	query.PriceMinimum = utils.StringToUint(c.Query("price_minimum"))
	query.PriceMaximum = utils.StringToUint(c.Query("price_maximum"))

	query.Page = utils.StringToUint(c.Query("page"))
	if query.Page == 0 {
		query.Page = 1
	}

	query.Limit = utils.StringToUint(c.Query("limit"))
	if query.Limit == 0 {
		query.Limit = 10
	}

	query.Sort = utils.QueryParamSort(c.Query("sort"))

	sortBy := strings.ToLower(c.Query("sort_by"))
	fmt.Println(sortBy)
	switch sortBy {
	case "newest":
		query.SortBy = "products.created_at"
	case "price":
		query.SortBy = "price_minimum"
	case "recommended":
		query.SortBy = "total_view_count"
	// case most buy
	default:
		query.SortBy = "name"
	}
	fmt.Println()

	return query
}
