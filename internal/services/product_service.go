package services

import (
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/product_dtos"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
)

type ProductService interface {
	GetAll(query *internalDtos.ProductQuery) ([]*internalDtos.Product, *utils.Pagination, error)
}

type productService struct {
	productRepository  r.ProductRepository
	merchantRepository r.MerchantRepository
}

type ProductServiceConfig struct {
	ProductRepository  r.ProductRepository
	MerchantRepository r.MerchantRepository
}

func NewProductService(c *ProductServiceConfig) ProductService {
	return &productService{
		productRepository:  c.ProductRepository,
		merchantRepository: c.MerchantRepository,
	}
}

func (s *productService) GetAll(query *internalDtos.ProductQuery) (productDetails []*internalDtos.Product, pagination *utils.Pagination, err error) {
	pagination = &utils.Pagination{
		Limit: query.Limit,
		Page:  query.Page,
	}

	productDetails, pagination, err = s.productRepository.GetAll(query, pagination)
	if err != nil {
		return
	}

	return productDetails, pagination, nil
}
