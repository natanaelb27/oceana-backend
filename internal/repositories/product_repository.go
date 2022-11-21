package repositories

import (
	"fmt"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos/product_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll(query *product_dtos.ProductQuery, pagination *utils.Pagination) ([]*product_dtos.Product, *utils.Pagination, error)
}

type productRepository struct {
	db *gorm.DB
}

type ProductRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductRepository(c *ProductRepositoryConfig) ProductRepository {
	return &productRepository{
		db: c.DB,
	}
}

func (r *productRepository) subQueryProductItemsMinimumMaximumPrice(priceType string) *gorm.DB {
	minMax := "MIN"
	if priceType == "max" {
		minMax = "MAX"
	}
	query := r.db.Model(&models.ProductItem{}).
		Select(minMax + "(product_items.price)").
		Where("product_items.product_id = products.id")
	return query
}

func (r *productRepository) subQueryProductViewsSum() *gorm.DB {
	query := r.db.Model(&models.ProductView{}).
		Select("SUM(product_views.count)").
		Where("product_views.product_id = products.id")
	return query
}

func (r *productRepository) subQueryProductCourierService(courierServices []string) *gorm.DB {
	query := r.db.Model(&models.ProductCourierService{}).
		Select("product_courier_services.courier_service_id").
		Where("product_courier_services.product_id = products.id").
		Where("product_courier_services.courier_service_id in (?)", courierServices)
	return query
}

func (r *productRepository) GetAll(queryParams *product_dtos.ProductQuery, pagination *utils.Pagination) ([]*product_dtos.Product, *utils.Pagination, error) {
	var ProductDetails []*product_dtos.Product

	query := r.db.Debug().Model(&models.Product{}).
		Select(`
			products.id,
            products.merchant_id,
            products.category_id,
            products.name,
            products.description,
            products.is_hazardous,
            products.weight,
            products.package_size,
            products.condition,
            products.product_status,
			(?) AS price_minimum,
			(?) AS price_maximum,
			(?) AS total_view_count,
            products.created_at,
            products.updated_at,
			merchants.id as "Merchant__id",
			merchants.name as "Merchant__name",
			cities.name as "Merchant__city"
		`,
			r.subQueryProductItemsMinimumMaximumPrice("min"),
			r.subQueryProductItemsMinimumMaximumPrice("max"),
			r.subQueryProductViewsSum(),
		).
		Joins("JOIN merchants ON merchants.id = products.merchant_id").
		Joins("JOIN addresses ON addresses.id = merchants.address_id").
		Joins("JOIN sub_districts ON sub_districts.id = addresses.sub_district_id").
		Joins("JOIN districts ON districts.id = sub_districts.district_id").
		Joins("JOIN cities ON cities.id = districts.city_id")

	if queryParams.MerchantID != 0 {
		query.Where("products.merchant_id = ?", queryParams.MerchantID)
	}

	if queryParams.CategoryID != 0 {
		query.Where("products.category_id = ?", queryParams.CategoryID)
	}

	if queryParams.PriceMinimum != 0 {
		query.Where(`(?) >= ?`, r.subQueryProductItemsMinimumMaximumPrice("min"), queryParams.PriceMinimum)
		query.Where(`(?) >= ?`, r.subQueryProductItemsMinimumMaximumPrice("min"), queryParams.PriceMinimum)
	}

	if queryParams.PriceMaximum != 0 {
		query.Where(`(?) <= ?`, r.subQueryProductItemsMinimumMaximumPrice("max"), queryParams.PriceMaximum)
		query.Where(`(?) <= ?`, r.subQueryProductItemsMinimumMaximumPrice("max"), queryParams.PriceMaximum)
	}

	if len(queryParams.MerchantCities) > 0 {
		query = query.Where("cities.id in (?)", queryParams.MerchantCities)
	}

	if len(queryParams.CourierServices) > 0 {
		fmt.Println(queryParams.CourierServices)
		query.Where(`EXISTS (?)`, r.subQueryProductCourierService(queryParams.CourierServices))
	}

	query.Where("products.name ILIKE ?", "%"+queryParams.Search+"%")

	query.
		Order(fmt.Sprintf("%s %s", queryParams.SortBy, queryParams.Sort)).
		Scopes(utils.Paginate(pagination, query)).
		Find(&ProductDetails)

	if query.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return ProductDetails, pagination, nil
}
