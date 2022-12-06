package services

import (
	"time"

	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/product_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type ProductService interface {
	GetAll(query *internalDtos.ProductQuery) ([]*internalDtos.Product, *utils.Pagination, error)
	ProductDetailsGetByID(productID uint) (*dtos.ProductDetailsGetByIDResponse, error)
	GetPurchaseQuantitiesByID(userID uuid.UUID, quantity *dtos.PurchaseQuantitiesQuery) (*dtos.GetPurchaseQuantitiesByIDResponse, error)
	ViewProduct(userID uuid.UUID, productID uint) (*dtos.GetProductViewResponse, error)
	GetFavoriteProduct(userID uuid.UUID, productID uint) (*dtos.GetProductFavoriteResponse, error)
	FavoriteProduct(userID uuid.UUID, productID uint) (*dtos.GetProductFavoriteResponse, error)
	GetUpdateByMerchantID(ID uint) (*internalDtos.ProductModel, error)
}

type productService struct {
	productRepository                r.ProductRepository
	merchantRepository               r.MerchantRepository
	productItemRepository            r.ProductItemRepository
	productVariationRepository       r.ProductVariationRepository
	productViewRepository            r.ProductViewRepository
	productFavoriteRepository        r.ProductFavoriteRepository
	merchantCourierServiceRepository r.MerchantCourierServiceRepository
}

type ProductServiceConfig struct {
	ProductRepository                r.ProductRepository
	MerchantRepository               r.MerchantRepository
	ProductItemRepository            r.ProductItemRepository
	ProductVariationRepository       r.ProductVariationRepository
	ProductViewRepository            r.ProductViewRepository
	ProductFavoriteRepository        r.ProductFavoriteRepository
	MerchantCourierServiceRepository r.MerchantCourierServiceRepository
}

func NewProductService(c *ProductServiceConfig) ProductService {
	return &productService{
		productRepository:                c.ProductRepository,
		merchantRepository:               c.MerchantRepository,
		productItemRepository:            c.ProductItemRepository,
		productVariationRepository:       c.ProductVariationRepository,
		productViewRepository:            c.ProductViewRepository,
		productFavoriteRepository:        c.ProductFavoriteRepository,
		merchantCourierServiceRepository: c.MerchantCourierServiceRepository,
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

	productVariation, err := s.productVariationRepository.GetByProductID(productID)

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

func (s *productService) ViewProduct(userID uuid.UUID, productID uint) (*dtos.GetProductViewResponse, error) {
	viewProduct, err := s.productViewRepository.GetByUserIDAndProductID(userID, productID)

	if err != nil && err != e.ErrDataNotFound {
		return nil, err
	}

	if err == e.ErrDataNotFound {
		model := &models.ProductView{
			UserId:       userID,
			ProductId:    productID,
			LastViewedAt: time.Now(),
			Count:        1,
		}

		newViewProduct, err := s.productViewRepository.Create(model)

		if err != nil {
			return nil, err
		}

		return dtos.NewGetProductViewResponse(newViewProduct), nil
	}

	if utils.IsDateMoreThanOneDay(viewProduct.LastViewedAt) {
		return dtos.NewGetProductViewResponse(viewProduct), nil
	}

	viewProduct.LastViewedAt = time.Now()
	viewProduct.Count += 1

	updatedViewProduct, err := s.productViewRepository.Update(viewProduct.ID, viewProduct)

	if err != nil {
		return nil, err
	}

	return dtos.NewGetProductViewResponse(updatedViewProduct), nil
}

func (s *productService) GetFavoriteProduct(userID uuid.UUID, productID uint) (*dtos.GetProductFavoriteResponse, error) {
	favoriteProduct, err := s.productFavoriteRepository.GetByUserIDAndProductID(userID, productID)

	if err != nil && err != e.ErrDataNotFound {
		return nil, err
	}

	if err == e.ErrDataNotFound {
		model := &models.ProductFavorite{
			UserId:     userID,
			ProductId:  productID,
			IsFavorite: false,
		}

		newFavoriteProduct, err := s.productFavoriteRepository.Create(model)

		if err != nil {
			return nil, err
		}

		return dtos.NewGetProductFavoriteResponse(newFavoriteProduct), nil
	}

	return dtos.NewGetProductFavoriteResponse(favoriteProduct), nil
}

func (s *productService) FavoriteProduct(userID uuid.UUID, productID uint) (*dtos.GetProductFavoriteResponse, error) {
	favoriteProduct, err := s.productFavoriteRepository.GetByUserIDAndProductID(userID, productID)

	if err != nil && err != e.ErrDataNotFound {
		return nil, err
	}

	if err == e.ErrDataNotFound {
		model := &models.ProductFavorite{
			UserId:     userID,
			ProductId:  productID,
			IsFavorite: false,
		}

		newFavoriteProduct, err := s.productFavoriteRepository.Create(model)

		if err != nil {
			return nil, err
		}

		return dtos.NewGetProductFavoriteResponse(newFavoriteProduct), nil
	}

	favoriteProduct.IsFavorite = !favoriteProduct.IsFavorite

	updatedFavoriteProduct, err := s.productFavoriteRepository.Update(favoriteProduct.ID, favoriteProduct)

	if err != nil {
		return nil, err
	}

	return dtos.NewGetProductFavoriteResponse(updatedFavoriteProduct), nil
}

func (s *productService) GetUpdateByMerchantID(merchantID uint) (*internalDtos.ProductModel, error) {
	// product
	productModel, err := s.productRepository.GetUpdateByMerchantID(merchantID)
	if err != nil {
		return nil, err
	}

	productItems, err := s.productItemRepository.GetUpdateByProductID(productModel.ID)
	// product_items
	if err != nil {
		return nil, err
	}
	productModel.ProductItems = productItems

	// variants
	variants, err := s.productVariationRepository.GetUpdateByProductID(productModel.ID)
	if err != nil {
		return nil, err
	}
	productModel.Variants = variants

	// get all courier on merchant courier service join to courier group by courier id
	// get all merchant courier service left join product courier service and add to its courier

	// couriers
	couriers, err := s.merchantCourierServiceRepository.GetAllCouriersByMerchantIDForProductUpdate(merchantID)
	if err != nil && err != e.ErrDataNotFound {
		return nil, err
	}
	productModel.Couriers = couriers
	//  courier > courier_services
	for iC := range couriers {
		courierServices, err := s.merchantCourierServiceRepository.GetAllCourierServicesByMerchantIDProductIDForProductUpdate(merchantID, productModel.ID)
		if err != nil && err != e.ErrDataNotFound {
			return nil, err
		}

		csIsActive := false
		for _, cs := range courierServices {
			if cs.IsActive {
				csIsActive = true
			}
		}

		productModel.Couriers[iC].CourierService = courierServices
		productModel.Couriers[iC].IsActive = csIsActive
	}

	return productModel, nil
}
