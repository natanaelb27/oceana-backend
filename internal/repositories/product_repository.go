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
			(
				SELECT MIN(product_items.price) 
				FROM product_items 
				WHERE product_items.product_id = products.id
			) AS price_minimum,
			(
				SELECT MAX(product_items.price) 
				FROM product_items 
				WHERE product_items.product_id = products.id
			) AS price_maximum,
			(
				SELECT SUM(product_views.count)
				FROM product_views 
				WHERE product_views.product_id = products.id
			) AS total_view_count,
            products.created_at,
            products.updated_at
		`)

	if queryParams.MerchantID != 0 {
		query.Where("products.merchant_id = ?", queryParams.MerchantID)
	}

	if queryParams.CategoryID != 0 {
		query.Where("products.category_id = ?", queryParams.CategoryID)
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
