package services

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants/product_media_constants"
	productInternalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/product_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type ProductService interface {
	GetAll(query *productInternalDtos.ProductQuery) ([]*productInternalDtos.Product, *utils.Pagination, error)
	ProductDetailsGetByID(productID uint) (*dtos.ProductDetailsGetByIDResponse, error)
	GetPurchaseQuantitiesByID(userID uuid.UUID, quantity *dtos.PurchaseQuantitiesQuery) (*dtos.GetPurchaseQuantitiesByIDResponse, error)
	ViewProduct(userID uuid.UUID, productID uint) (*dtos.GetProductViewResponse, error)
	GetFavoriteProduct(userID uuid.UUID, productID uint) (*dtos.GetProductFavoriteResponse, error)
	FavoriteProduct(userID uuid.UUID, productID uint) (*dtos.GetProductFavoriteResponse, error)
	CreateProduct(request *productInternalDtos.ProductModel) (*productInternalDtos.ProductModel, error)
	GetUpdateModelByProductID(productID uint) (*productInternalDtos.ProductModel, error)
}

type productService struct {
	productRepository                r.ProductRepository
	merchantRepository               r.MerchantRepository
	productItemRepository            r.ProductItemRepository
	productReviewService             r.ProductReviewRepository
	productVariationRepository       r.ProductVariationRepository
	productViewRepository            r.ProductViewRepository
	productFavoriteRepository        r.ProductFavoriteRepository
	merchantCourierServiceRepository r.MerchantCourierServiceRepository
	categoryRepository               r.CategoryRepository
	productVariantTypeRepository     r.ProductVariantTypeRepository
	productVariantMediaRepository    r.ProductVariantMediaRepository
	productVariantRepository         r.ProductVariantRepository
	productCourierServiceRepository  r.ProductCourierServiceRepository
	productMediaRepository           r.ProductMediaRepository
}

type ProductServiceConfig struct {
	ProductRepository                r.ProductRepository
	MerchantRepository               r.MerchantRepository
	ProductItemRepository            r.ProductItemRepository
	ProductReviewRepository          r.ProductReviewRepository
	ProductVariationRepository       r.ProductVariationRepository
	ProductViewRepository            r.ProductViewRepository
	ProductFavoriteRepository        r.ProductFavoriteRepository
	MerchantCourierServiceRepository r.MerchantCourierServiceRepository
	CategoryRepository               r.CategoryRepository
	ProductVariantTypeRepository     r.ProductVariantTypeRepository
	ProductVariantMediaRepository    r.ProductVariantMediaRepository
	ProductVariantRepository         r.ProductVariantRepository
	ProductCourierServiceRepository  r.ProductCourierServiceRepository
	ProductMediaRepository           r.ProductMediaRepository
}

func NewProductService(c *ProductServiceConfig) ProductService {
	return &productService{
		productRepository:                c.ProductRepository,
		merchantRepository:               c.MerchantRepository,
		productItemRepository:            c.ProductItemRepository,
		productReviewService:             c.ProductReviewRepository,
		productVariationRepository:       c.ProductVariationRepository,
		productViewRepository:            c.ProductViewRepository,
		productFavoriteRepository:        c.ProductFavoriteRepository,
		merchantCourierServiceRepository: c.MerchantCourierServiceRepository,
		categoryRepository:               c.CategoryRepository,
		productVariantTypeRepository:     c.ProductVariantTypeRepository,
		productVariantMediaRepository:    c.ProductVariantMediaRepository,
		productVariantRepository:         c.ProductVariantRepository,
		productCourierServiceRepository:  c.ProductCourierServiceRepository,
		productMediaRepository:           c.ProductMediaRepository,
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

	productReviewDetails, err := s.productReviewService.GetDetailsByProductID(productID)

	if err != nil {
		return nil, err
	}

	productDetails.ProductReviewDetails = productReviewDetails

	return dtos.NewProductDetailsGetByIDResponse(productDetails), nil
}

func (s *productService) GetPurchaseQuantitiesByID(userID uuid.UUID, query *dtos.PurchaseQuantitiesQuery) (*dtos.GetPurchaseQuantitiesByIDResponse, error) {
	purchaseQuantities, err := s.productRepository.GetPurchaseQuantitiesByID(userID, query.ProductID, query.ProductItemID, query.Quantity)

	if err != nil {
		return nil, err
	}

	return dtos.NewGetPurchaseQuantitiesByIDResponse(purchaseQuantities), nil
}

func (s *productService) GetAll(query *productInternalDtos.ProductQuery) (productDetails []*productInternalDtos.Product, pagination *utils.Pagination, err error) {
	pagination = &utils.Pagination{
		Limit: query.Limit,
		Page:  query.Page,
	}

	if query.CategoryID != 0 {
		var categories []uint
		categories, err = s.GetCategoryIds(query.CategoryID)
		if err != nil {
			return
		}

		query.Categories = categories
	}

	productDetails, pagination, err = s.productRepository.GetAll(query, pagination)
	if err != nil {
		return
	}

	return productDetails, pagination, nil
}

func (d *productService) GetAllSubCategories(categoryData []*models.Category) ([]uint, error) {
	category := []uint{}

	for _, cat := range categoryData {
		subCategoryData, err := d.categoryRepository.GetAllSubCategoriesByParentID(cat.ID)
		if err != nil {
			return nil, err
		}

		category = append(category, cat.ID)

		if len(subCategoryData) > 0 {
			subCategory, err := d.GetAllSubCategories(subCategoryData)
			if err != nil {
				return nil, err
			}

			category = append(category, subCategory...)
		}
	}

	return category, nil
}

func (s *productService) GetCategoryIds(id uint) ([]uint, error) {
	categoryData, err := s.categoryRepository.GetAllSubCategoriesByParentID(id)
	if err != nil {
		return nil, err
	}

	category, err := s.GetAllSubCategories(categoryData)
	if err != nil {
		return nil, err
	}

	category = append(category, id)

	return category, nil
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

func (s *productService) CreateProduct(request *productInternalDtos.ProductModel) (*productInternalDtos.ProductModel, error) {
	// create product
	// check if category exist
	product := &models.Product{
		MerchantID:  request.MerchantID,
		CategoryID:  request.CategoryID,
		Name:        request.Name,
		Description: request.Description,
		IsHazardous: request.IsHazardous,
		Weight:      request.Weight,
		PackageSize: request.PackageSize,
		Condition:   request.Condition,
		Status:      request.Status,
		ProcessDay:  request.ProcessDay,
		PreorderDay: request.PreorderDay,
	}
	createdProduct, err := s.productRepository.CreateProduct(product)
	if err != nil {
		return nil, err
	}

	// create product media from images
	for _, iI := range request.Images {
		product := &models.ProductMedia{
			ProductID: product.ID,
			MediaURL:  iI,
			Type:      product_media_constants.TypeImage,
		}
		_, err := s.productMediaRepository.Create(product)
		if err != nil {
			return nil, err
		}
	}

	// create product media from videos
	for _, iV := range request.Videos {
		product := &models.ProductMedia{
			ProductID: product.ID,
			MediaURL:  iV,
			Type:      product_media_constants.TypeVideo,
		}
		_, err := s.productMediaRepository.Create(product)
		if err != nil {
			return nil, err
		}
	}

	// map product variant for later use to identify product variations
	productVariantMap := make([][]uint, len(request.Variants))
	for iVariant, variant := range request.Variants {
		// declare product variant map with options length
		productVariantMap[iVariant] = make([]uint, len(variant.Options))
		// create product variant types
		productVariantType := &models.ProductVariantType{
			Name:      variant.Name,
			ProductID: createdProduct.ID,
		}
		createdProductVariantType, err := s.productVariantTypeRepository.Create(productVariantType)
		if err != nil {
			return nil, err
		}

		for iOption, option := range variant.Options {
			// create product variant
			productVariants := &models.ProductVariant{
				Name:                 option,
				ProductVariantTypeID: createdProductVariantType.ID,
			}
			createdProductVariant, err := s.productVariantRepository.Create(productVariants)
			if err != nil {
				return nil, err
			}

			// set product variant id to product variant map with certain position
			productVariantMap[iVariant][iOption] = productVariants.ID

			// create product variant media on first variant only
			if iVariant == 0 {
				productVariantMedia := &models.ProductVariantMedia{
					ProductVariantID: createdProductVariant.ID,
					MediaURL:         variant.Images[iOption],
				}
				_, err = s.productVariantMediaRepository.Create(productVariantMedia)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// create product item
	for _, pi := range request.ProductItems {
		productItem := &models.ProductItem{
			ProductID:     product.ID,
			SKU:           pi.SKU,
			Price:         pi.Price,
			StockQuantity: pi.StockQuantity,
		}
		createdProductItem, err := s.productItemRepository.CreateProductItem(productItem)
		if err != nil {
			return nil, err
		}
		// create product variation
		for iVi, vi := range pi.VariantIndex {
			productVariation := &models.ProductVariation{
				ProductItemID:    createdProductItem.ID,
				ProductVariantID: productVariantMap[iVi][vi],
			}
			_, err := s.productVariationRepository.Create(productVariation)
			if err != nil {
				return nil, err
			}

		}
	}

	for _, c := range request.Couriers {
		for _, cs := range c.CourierServices {
			// if is active false then continue to next iteration
			if !cs.IsActive {
				continue
			}
			// get merchant courier service
			merchantCourierServiceConditionMap := map[string]any{
				"courier_service_id": cs.ID,
				"merchant_id":        request.MerchantID,
				"is_active":          true,
			}
			merchantCourierService, err := s.merchantCourierServiceRepository.GetByConditionMap(merchantCourierServiceConditionMap)
			if err != nil {
				return nil, err
			}
			// create product courier service
			productCourierService := &models.ProductCourierService{
				MerchantCourierServiceID: merchantCourierService.ID,
				ProductId:                createdProduct.ID,
			}
			_, err = s.productCourierServiceRepository.Create(productCourierService)
			if err != nil {
				return nil, err
			}
		}
	}

	// get created product model
	productModel, err := s.GetUpdateModelByProductID(createdProduct.ID)
	if err != nil {
		return nil, err
	}

	return productModel, nil
}

func (s *productService) GetUpdateModelByProductID(productID uint) (*productInternalDtos.ProductModel, error) {
	// product
	productModel, err := s.productRepository.GetUpdateByProductID(productID)
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
	couriers, err := s.merchantCourierServiceRepository.GetAllCouriersByMerchantIDForProductUpdate(productModel.MerchantID)
	if err != nil && err != e.ErrDataNotFound {
		return nil, err
	}
	productModel.Couriers = couriers

	//  courier > courier_services
	for iC := range couriers {
		courierServices, err := s.merchantCourierServiceRepository.GetAllCourierServicesByMerchantIDProductIDForProductUpdate(productModel.MerchantID, productModel.ID)
		if err != nil && err != e.ErrDataNotFound {
			return nil, err
		}

		csIsActive := false
		for _, cs := range courierServices {
			if cs.IsActive {
				csIsActive = true
			}
		}

		productModel.Couriers[iC].CourierServices = courierServices
		productModel.Couriers[iC].IsActive = csIsActive
	}

	return productModel, nil
}
