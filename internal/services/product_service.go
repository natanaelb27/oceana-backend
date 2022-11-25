package services

import (
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/product_dtos"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/google/uuid"
)

type ProductService interface {
	GetAll(query *internalDtos.ProductQuery) ([]*internalDtos.Product, *utils.Pagination, error)
	ProductDetailsGetByID(productID uint) (*dtos.ProductDetailsGetByIDResponse, error)
	GetPurchaseQuantitiesByID(userID uuid.UUID, quantity *dtos.PurchaseQuantitiesQuery) (*dtos.GetPurchaseQuantitiesByIDResponse, error)
}

type productService struct {
	productRepository        r.ProductRepository
	merchantRepository       r.MerchantRepository
	productItemRepository    r.ProductItemRepository
	productVariantRepository r.ProductVariationRepository
}

type ProductServiceConfig struct {
	ProductRepository          r.ProductRepository
	MerchantRepository         r.MerchantRepository
	ProductItemRepository      r.ProductItemRepository
	ProductVariationRepository r.ProductVariationRepository
}

func NewProductService(c *ProductServiceConfig) ProductService {
	return &productService{
		productRepository:        c.ProductRepository,
		merchantRepository:       c.MerchantRepository,
		productItemRepository:    c.ProductItemRepository,
		productVariantRepository: c.ProductVariationRepository,
	}
}

func (s *productService) ProductDetailsGetByID(productID uint) (*dtos.ProductDetailsGetByIDResponse, error) {
	productDetails, err := s.productRepository.ProductDetailsGetByID(productID)

	if err != nil {
		return nil, err
	}

	productModels, err := s.productItemRepository.GetByProductID(productID)

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

func (s *productService) GetPurchaseQuantitiesByID(userID uuid.UUID, query *dtos.PurchaseQuantitiesQuery) (*dtos.GetPurchaseQuantitiesByIDResponse, error) {
	purchaseQuantities, err := s.productRepository.GetPurchaseQuantitiesByID(userID, query.ProductID, query.ProductItemID, query.Quantity)

	if err != nil {
		return nil, err
	}

	return dtos.NewGetPurchaseQuantitiesByIDResponse(purchaseQuantities), nil
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
