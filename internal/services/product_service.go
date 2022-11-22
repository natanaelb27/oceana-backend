package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
)

type ProductService interface {
	CreateProduct(product *models.Product) (*models.Product, error)
}

type productService struct {
	productRepository r.ProductRepository
}

type ProductServiceConfig struct {
	ProductRepository r.ProductRepository
}

func NewProductService(c *ProductServiceConfig) ProductService {
	return &productService{
		productRepository: c.ProductRepository,
	}
}

func (p *productService) CreateProduct(model *models.Product) (*models.Product, error) {
	return p.productRepository.CreateProduct(model)
}
