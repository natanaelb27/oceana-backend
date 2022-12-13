package product_dtos

import (
	"strconv"
	"strings"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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
	Name                          string    `json:"name"`
	Description                   string    `json:"description"`
	MediaURL                      string    `json:"media_url"`
	IsHazardous                   bool      `json:"is_hazardous"`
	Weight                        uint      `json:"weight"`
	PackageSize                   string    `json:"package_size"`
	Condition                     string    `json:"condition"`
	ProductStatus                 string    `json:"product_status"`
	PriceMinimum                  uint      `json:"price_minimum"`
	PriceMaximum                  uint      `json:"price_maximum"`
	PriceMinimumAfterDiscount     float64   `json:"price_minimum_after_discount"`
	PriceMaximumAfterDiscount     float64   `json:"price_maximum_after_discount"`
	ProductMerchantPromotionID    uint      `json:"product_merchant_promotion_id"`
	MerchantPromotionName         string    `json:"merchant_promotion_name"`
	Discount                      uint      `json:"discount"`
	DiscountType                  string    `json:"discount_type"`
	DiscountMaximumDiscountAmount uint      `json:"discount_maximum_discount_amount"`
	TotalViewCount                uint      `json:"total_view_count"`
	TotalPurchaseCount            uint      `json:"total_purchase_count"`
	RatingStar                    float64   `json:"rating_star"`
	Merchant                      Merchant  `json:"merchant"`
	CreatedAt                     time.Time `json:"created_at"`
	UpdatedAt                     time.Time `json:"updated_at"`
}

type ProductQuery struct {
	Search string
	SortBy string
	Sort   string
	Limit  uint
	Page   uint
	// Filter
	MerchantID        uint
	CategoryID        uint
	PriceMinimum      uint
	PriceMaximum      uint
	RatingStarMinimum uint
	MerchantCities    []string
	CourierServices   []string
	// Categories
	Categories []uint
}

func NewProductQuery(c *gin.Context) *ProductQuery {
	query := &ProductQuery{}

	query.Search = c.Query("search")
	query.MerchantID = utils.StringToUint(c.Query("merchant-id"))
	query.CategoryID = utils.StringToUint(c.Query("category-id"))
	query.PriceMinimum = utils.StringToUint(c.Query("price-minimum"))
	query.PriceMaximum = utils.StringToUint(c.Query("price-maximum"))
	query.RatingStarMinimum = utils.StringToUint(c.Query("rating-star-minimum"))

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
		query.SortBy = "price_minimum_after_discount"
	case "most-buy":
		query.SortBy = "total_purchase_count"
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

type CourierService struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Description string `json:"name"`
	IsActive    bool   `json:"is_active"`
}

type Courier struct {
	ID              uint              `json:"id"`
	Code            string            `json:"code"`
	Name            string            `json:"name"`
	IsActive        bool              `json:"is_active" gorm:"-"`
	CourierServices []*CourierService `json:"courier_services" gorm:"-"`
}

type Logistics struct {
	Price     uint `json:"price"`
	IsEnabled bool `json:"is_enabled"`
}

type Variants struct {
	Name    string         `json:"name"`
	Options pq.StringArray `json:"options" gorm:"type:text[]"`
	Images  pq.StringArray `json:"images" gorm:"type:text[]"`
}

type ProductItem struct {
	ID  uint   `json:"id"`
	SKU string `json:"sku"`
	// pq.Int64Array which is integer[] in gorm can receive array of integer from sql (for example array_agg)
	VariantIndex  pq.Int64Array `json:"variant_index" gorm:"type:integer[]"`
	Price         uint          `json:"price"`
	StockQuantity uint          `json:"stock_quantity"`
}

type ProductModel struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CategoryID  uint           `json:"category_id"`
	MerchantID  uint           `json:"merchant_id"`
	Weight      uint           `json:"weight"`
	IsHazardous bool           `json:"is_hazardous"`
	PackageSize string         `json:"package_size"`
	Condition   string         `json:"condition"`
	Status      string         `json:"status"`
	ProcessDay  uint           `json:"process_day"`
	PreorderDay uint           `json:"preorder_day"`
	Images      pq.StringArray `json:"images" gorm:"type:text[]"`
	Videos      pq.StringArray `json:"videos" gorm:"type:text[]"`
	// gorm:"-" tells gorm ignore this field when write and read with struct
	ProductItems []*ProductItem `json:"product_items" gorm:"-"`
	Variants     []*Variants    `json:"variants" binding:"min=1" gorm:"-"`
	Couriers     []*Courier     `json:"couriers" binding:"min=1" gorm:"-"`
}

// product variant type update info
type PvtUpdateInfo struct {
	ID     uint
	NewPVT *models.ProductVariantType
}

// product variant update info
type PvUpdateInfo struct {
	Before *models.ProductVariant
	After  *models.ProductVariant
}

// product variant media update info
type PvmUpdateInfo struct {
	Before *models.ProductVariantMedia
	After  *models.ProductVariantMedia
}

// product items update info
type PiUpdateInfo struct {
	Before *models.ProductItem
	After  *models.ProductItem
}

// courier services update info
type PcsUpdateInfo struct {
	Before *models.ProductCourierService
	After  *models.ProductCourierService
}
