package services

import (
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"

	"github.com/google/uuid"
)

type ProductService interface {
	ProductDetailsGetByID(userID uuid.UUID, productID uint) (*dtos.ProductDetailsGetByIDResponse, error)
}

type productService struct {
	productRepository        r.ProductRepository
	productItemRepository    r.ProductItemRepository
	productVariantRepository r.ProductVariationRepository
}

type ProductServiceConfig struct {
	ProductRepository          r.ProductRepository
	ProductItemRepository      r.ProductItemRepository
	ProductVariationRepository r.ProductVariationRepository
}

func NewProductService(c *ProductServiceConfig) ProductService {
	return &productService{
		productRepository:        c.ProductRepository,
		productItemRepository:    c.ProductItemRepository,
		productVariantRepository: c.ProductVariationRepository,
	}
}
func (s *productService) ProductDetailsGetByID(userID uuid.UUID, productID uint) (*dtos.ProductDetailsGetByIDResponse, error) {
	productDetails, err := s.productRepository.ProductDetailsGetByID(userID, productID)

	if err != nil {
		return nil, err
	}

	productModels, err := s.productItemRepository.GetByProductID(userID, productID)

	if err != nil {
		return nil, err
	}

	productDetails.ProductModel = productModels

	productVariation, err := s.productVariantRepository.GetByProductID(productID)

	if err != nil {
		return nil, err
	}

	productDetails.ProductVariation = productVariation

	return dtos.NewProductDetailsGetByIDResponse(productDetails), nil
}
