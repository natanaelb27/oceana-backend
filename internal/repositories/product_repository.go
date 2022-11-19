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

	query := r.db.Model(&models.Product{}).
		Select(`*`)

	if queryParams.MerchantID != 0 {
		query.Where("merchant_id = ?", queryParams.MerchantID)
	}

	query.Where("name ILIKE ?", "%"+queryParams.Search+"%")

	query.
		Order(fmt.Sprintf("%s %s", queryParams.SortBy, queryParams.Sort)).
		Scopes(utils.Paginate(pagination, query)).
		Find(&ProductDetails)

	if query.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return ProductDetails, pagination, nil
}
