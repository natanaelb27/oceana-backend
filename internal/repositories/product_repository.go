package repositories

import (
	"database/sql"
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
	sqPriceMinimum := sql.Named("price_minimum", r.subQueryProductItemsMinimumMaximumPrice("min"))
	sqPriceMaximum := sql.Named("price_maximum", r.subQueryProductItemsMinimumMaximumPrice("max"))
	sqTotalViewCount := sql.Named("total_view_count", r.subQueryProductViewsSum())
	// subquery promotion price

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
			(@price_minimum) AS price_minimum,
			(@price_maximum) AS price_maximum,
			-- count price minimum after discount
			(
				CASE
					WHEN pmp.discount_type = 'percentage' 
						THEN (@price_minimum) - LEAST
							(
								(pmp.discount / 100.0) * (@price_minimum), 
								pmp.maximum_discount_amount
							)
					WHEN pmp.discount_type = 'nominal' 
						THEN GREATEST
							(
								(@price_minimum) - pmp.discount, 
								0
							)
				ELSE (@price_minimum)
				END
			) AS price_minimum_after_discount,
			-- count price maximum after discount
			(
				CASE
					WHEN pmp.discount_type = 'percentage' 
						THEN (@price_maximum) - LEAST
							(
								(pmp.discount / 100.0) * (@price_maximum), 
								pmp.maximum_discount_amount
							)
					WHEN pmp.discount_type = 'nominal' 
						THEN GREATEST
							(
								(@price_maximum) - pmp.discount, 
								0
							)
				ELSE (@price_maximum)
				END
			) AS price_maximum_after_discount,
			(@total_view_count) AS total_view_count,
			pmp.discount as discount,
            products.created_at,
            products.updated_at,
			merchants.id as "Merchant__id",
			merchants.name as "Merchant__name",
			cities.name as "Merchant__city"
		`,
			sqPriceMinimum,
			sqPriceMaximum,
			sqTotalViewCount,
		).
		Joins("JOIN merchants ON merchants.id = products.merchant_id").
		Joins("JOIN addresses ON addresses.id = merchants.address_id").
		Joins("JOIN sub_districts ON sub_districts.id = addresses.sub_district_id").
		Joins("JOIN districts ON districts.id = sub_districts.district_id").
		Joins("JOIN cities ON cities.id = districts.city_id").
		// left join to get merchant promotion
		Joins("LEFT JOIN product_merchant_promotions pmp ON products.id = pmp.product_id")

	// group when status is active or when product got no promotion
	query.Group("products.id, pmp.id, merchants.id, cities.id").Having("pmp.status = 'active' OR pmp.id IS NULL")

	// filter by merchant
	if queryParams.MerchantID != 0 {
		query.Where("products.merchant_id = ?", queryParams.MerchantID)
	}

	// filter by price category
	if queryParams.CategoryID != 0 {
		query.Where("products.category_id = ?", queryParams.CategoryID)
	}

	// filter by price minimum
	if queryParams.PriceMinimum != 0 {
		query.Where(`(@price_minimum) >= @query_price_minimum`,
			sqPriceMinimum,
			sql.Named("query_price_minimum", queryParams.PriceMinimum),
		)
	}
	// filter by price maximum
	if queryParams.PriceMaximum != 0 {
		query.Where(`(@price_maximum) <= @query_price_maximum`,
			sqPriceMaximum,
			sql.Named("query_price_maximum", queryParams.PriceMaximum),
		)
	}

	// filter by merchant cities
	if len(queryParams.MerchantCities) > 0 {
		query = query.Where("cities.id in (?)", queryParams.MerchantCities)
	}

	// filter by cities
	if len(queryParams.CourierServices) > 0 {
		query.Where(`EXISTS (?)`,
			r.subQueryProductCourierService(queryParams.CourierServices),
		)
	}

	// search by products name
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
