package services

import (
	"fmt"
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/constants/product_media_constants"
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	productInternalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/product_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService interface {
	WithTx(tx *gorm.DB) *productService
	GetAll(query *productInternalDtos.ProductQuery) ([]*productInternalDtos.Product, *utils.Pagination, error)
	ProductDetailsGetByID(internalQuery *internalDtos.ProductDetailsInternalQuery) (*dtos.ProductDetailsGetByIDResponse, error)
	GetPurchaseQuantitiesByID(userID uuid.UUID, quantity *dtos.PurchaseQuantitiesQuery) (*dtos.GetPurchaseQuantitiesByIDResponse, error)
	ViewProduct(userID uuid.UUID, productID uint) (*dtos.GetProductViewResponse, error)
	GetFavoriteProduct(userID uuid.UUID, productID uint) (*dtos.GetProductFavoriteResponse, error)
	FavoriteProduct(userID uuid.UUID, productID uint) (*dtos.GetProductFavoriteResponse, error)
	CreateProduct(request *productInternalDtos.ProductModel, userID uuid.UUID) (*productInternalDtos.ProductModel, error)
	UpdateProductUsingModel(userID uuid.UUID, request *productInternalDtos.ProductModel) (*productInternalDtos.ProductModel, error)
	GetUpdateModelByProductID(productID uint) (*productInternalDtos.ProductModel, error)
	GetAvailableProductsByPromotionPeriod(id uuid.UUID, query *dtos.AvailableProductsByPromotionPeriodQuery) ([]*dtos.AvailableProductsResponse, error)
	CheckCourierService(request *productInternalDtos.ProductCheckCourierService) (*productInternalDtos.ProductCheckCourierServiceResponse, error)
	MyShopGetAll(query *productInternalDtos.ProductQuery) (responseProductDetails []*dtos.ProductDetailsGetByIDResponse, pagination *utils.Pagination, err error)
	GetAllUserFavoriteProducts(userID uuid.UUID, query *productInternalDtos.ProductFavoriteQuery) (productDetails []*productInternalDtos.Product, pagination *utils.Pagination, err error)
	DeleteFavoriteProducts(userID uuid.UUID, request *dtos.DeleteFavoriteProductsRequest) ([]*dtos.UpdateProductFavoriteResponse, error)
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

func (u productService) WithTx(tx *gorm.DB) *productService {
	u.productRepository = u.productRepository.WithTx(tx)
	u.merchantRepository = u.merchantRepository.WithTx(tx)
	u.productItemRepository = u.productItemRepository.WithTx(tx)
	u.productVariationRepository = u.productVariationRepository.WithTx(tx)
	u.productViewRepository = u.productViewRepository.WithTx(tx)
	u.productFavoriteRepository = u.productFavoriteRepository.WithTx(tx)
	u.merchantCourierServiceRepository = u.merchantCourierServiceRepository.WithTx(tx)
	u.categoryRepository = u.categoryRepository.WithTx(tx)
	u.productVariantTypeRepository = u.productVariantTypeRepository.WithTx(tx)
	u.productVariantMediaRepository = u.productVariantMediaRepository.WithTx(tx)
	u.productVariantRepository = u.productVariantRepository.WithTx(tx)
	u.productCourierServiceRepository = u.productCourierServiceRepository.WithTx(tx)
	u.productMediaRepository = u.productMediaRepository.WithTx(tx)

	return &u
}

func (s *productService) ProductDetailsGetByID(internalQuery *internalDtos.ProductDetailsInternalQuery) (*dtos.ProductDetailsGetByIDResponse, error) {
	productID := internalQuery.ProductID

	productDetails, err := s.productRepository.GetDetails(internalQuery)

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
			IsFavorite: true,
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

func (s *productService) CreateProduct(request *productInternalDtos.ProductModel, userID uuid.UUID) (*productInternalDtos.ProductModel, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	request.MerchantID = merchant.ID

	if request.Status != constants.ProductStatusActive &&
		request.Status != constants.ProductStatusInactive {
		request.Status = constants.ProductStatusInactive
	}

	// create product
	// check if category exist
	product := &models.Product{
		MerchantID:  request.MerchantID,
		CategoryID:  request.CategoryID,
		Name:        request.Name,
		Description: request.Description,
		IsHazardous: *request.IsHazardous,
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
				if len(variant.Images) != len(variant.Options) {
					return nil, e.ErrEmptyVariants
				}

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
			if !cs.IsSelected {
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
	fmt.Println("S")
	// get created product model
	productModel, err := s.GetUpdateModelByProductID(createdProduct.ID)
	if err != nil {
		return nil, err
	}

	return productModel, nil
}

func (s *productService) UpdateProductUsingModel(userID uuid.UUID, request *productInternalDtos.ProductModel) (*productInternalDtos.ProductModel, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	// check if userID merchantID is = model.MerchantID
	if request.MerchantID != merchant.ID {
		return nil, e.ErrUnauthorized
	}

	product, err := s.productRepository.GetById(request.ID)
	if err != nil {
		return nil, err
	}
	// check if user owns model.id (product id)
	if product.MerchantID != merchant.ID {
		return nil, e.ErrUnauthorized
	}

	// get original product model
	oriProductModel, err := s.GetUpdateModelByProductID(product.ID)
	if err != nil {
		return nil, err
	}

	//* products
	// product validaiton
	// check if product status not its enum
	if request.Status != constants.ProductStatusActive &&
		request.Status != constants.ProductStatusInactive {
		return nil, e.ErrInvalidModelStatus
	}
	// check if product condition not its enum
	if request.Condition != constants.ProductConditionNew &&
		request.Condition != constants.ProductConditionUsed {
		return nil, e.ErrInvalidModel
	}
	newProductMap := map[string]any{
		"id":           request.ID,
		"name":         request.Name,
		"description":  request.Description,
		"category_id":  request.CategoryID,
		"merchant_id":  request.MerchantID,
		"weight":       request.Weight,
		"is_hazardous": request.IsHazardous,
		"package_size": request.PackageSize,
		"condition":    request.Condition,
		"status":       request.Status,
		"process_day":  request.ProcessDay,
		"preorder_day": request.PreorderDay,
	}
	// update product
	err = s.productRepository.UpdateUsingMap(product.ID, newProductMap)
	if err != nil {
		return nil, err
	}

	// product medias to create
	var pmCreates = make([]*productInternalDtos.PmUpdateInfo, 0)
	// product medias to delete
	var pmDeletes = make([]*productInternalDtos.PmUpdateInfo, 0)

	//* product medias - check changes from images
	// create product media from images
	// map images original and request
	// original images for product media map
	var oriIPMMap = map[string]bool{}
	for _, v := range oriProductModel.Images {
		oriIPMMap[v] = true
	}
	// request images for product media map
	var reqIPMMap = map[string]bool{}
	for _, v := range request.Images {
		reqIPMMap[v] = true
	}
	// add product media to create
	for _, v := range request.Images {
		if exists := oriIPMMap[v]; !exists {
			pmCreates = append(pmCreates, &productInternalDtos.PmUpdateInfo{
				After: &models.ProductMedia{
					ProductID: product.ID,
					MediaURL:  v,
					Type:      product_media_constants.TypeImage,
				},
			})
		}
	}
	// add product media to delete
	for _, v := range oriProductModel.Images {
		if exists := reqIPMMap[v]; !exists {
			pmDeletes = append(pmDeletes, &productInternalDtos.PmUpdateInfo{
				Before: &models.ProductMedia{
					ProductID: product.ID,
					MediaURL:  v,
					Type:      product_media_constants.TypeImage,
				},
			})
		}
	}

	//* product mdias - check changes from videos
	// create product media from videos
	// map videos original and request
	// original videos for product media map
	var oriVPMMap = map[string]bool{}
	for _, v := range oriProductModel.Videos {
		oriVPMMap[v] = true
	}
	// request images for product media map
	var reqVPMMap = map[string]bool{}
	for _, v := range request.Videos {
		reqVPMMap[v] = true
	}
	// add product media to create
	for _, v := range request.Videos {
		if exists := oriVPMMap[v]; !exists {
			pmCreates = append(pmCreates, &productInternalDtos.PmUpdateInfo{
				After: &models.ProductMedia{
					ProductID: product.ID,
					MediaURL:  v,
					Type:      product_media_constants.TypeVideo,
				},
			})
		}
	}
	// add product media to delete
	for _, v := range oriProductModel.Videos {
		if exists := reqVPMMap[v]; !exists {
			pmDeletes = append(pmDeletes, &productInternalDtos.PmUpdateInfo{
				Before: &models.ProductMedia{
					ProductID: product.ID,
					MediaURL:  v,
					Type:      product_media_constants.TypeVideo,
				},
			})
		}
	}

	piLengthValidation := 0

	if len(request.ProductItems) > 0 {
		piLengthValidation = 1
	}

	var pvtCreates = []*productInternalDtos.PvtUpdateInfo{}
	var pvtUpdates = []*productInternalDtos.PvtUpdateInfo{}
	var pvtDeletes = []*productInternalDtos.PvtUpdateInfo{}

	var pvCreates = []*productInternalDtos.PvUpdateInfo{}
	var pvUpdates = []*productInternalDtos.PvUpdateInfo{}
	var pvDeletes = []*productInternalDtos.PvUpdateInfo{}

	var pvmCreates = []*productInternalDtos.PvmUpdateInfo{}
	var pvmUpdates = []*productInternalDtos.PvmUpdateInfo{}
	var pvmDeletes = []*productInternalDtos.PvmUpdateInfo{}

	// get product variant type id [n] by getting product variant type by product id
	// current product variant types
	var pvts = []*models.ProductVariantType{}
	pvts, err = s.productVariantTypeRepository.GetAllByProductID(oriProductModel.ID)
	if err != nil && err != e.ErrDataNotFound {
		return nil, err
	}

	//* product variant types - check changes
	// product variant types - req to ori comparison - append creates
	if len(request.Variants) > len(oriProductModel.Variants) {
		for i := len(oriProductModel.Variants); i < len(request.Variants); i++ {
			pvtCreates = append(pvtCreates, &productInternalDtos.PvtUpdateInfo{
				ID: 0,
				NewPVT: &models.ProductVariantType{
					Name:      request.Variants[i].Name,
					ProductID: product.ID,
				},
			})
		}
	}
	// product variant types - ori to req comparison - append updates and deletes
	for iOV, oV := range oriProductModel.Variants {
		// if current array index exceed request variants length
		if iOV+1 > len(request.Variants) {
			pvtDeletes = append(pvtDeletes, &productInternalDtos.PvtUpdateInfo{
				ID: pvts[iOV].ID,
			})
			break
		}
		// request variants on original variants current index
		var rV = request.Variants[iOV]
		if oV.Name != rV.Name {
			pvtUpdates = append(pvtUpdates, &productInternalDtos.PvtUpdateInfo{
				ID: pvts[iOV].ID,
				NewPVT: &models.ProductVariantType{
					Name: rV.Name,
				},
			})
		}

		fmt.Println("MASUK")
		//* product variants - check changes
		// validation
		// if length of images not equal options then invalid model
		if len(rV.Options) <= 0 {
			return nil, e.ErrEmptyVariants
		}
		piLengthValidation *= len(rV.Options)

		// get all product variants from database by productVariantTypeID and Name
		pv, err := s.productVariantRepository.GetAllByProductVariantTypeID(pvts[iOV].ID)
		if err != nil {
			return nil, err
		}
		// product variants - req to ori comparison - append creates
		if len(rV.Options) > len(oV.Options) {
			for i := len(oV.Options); i < len(rV.Options); i++ {
				pvCreates = append(pvCreates, &productInternalDtos.PvUpdateInfo{
					After: &models.ProductVariant{
						Name:                 rV.Options[i],
						ProductVariantTypeID: pvts[iOV].ID,
					},
				})
			}
		}
		// product variants - ori to req comparison - append updates and deletes
		for iOO, oO := range oV.Options {
			rO := rV.Options
			// if index not exist then append deletes
			if len(rO)-1 < iOO {
				pvDeletes = append(pvDeletes, &productInternalDtos.PvUpdateInfo{
					Before: &models.ProductVariant{
						ID: pv[iOO].ID,
					},
				})
				continue
			} else if oO != rO[iOO] {
				// if there is changes then append updates
				pvUpdates = append(pvUpdates, &productInternalDtos.PvUpdateInfo{
					Before: &models.ProductVariant{
						ID: pv[iOO].ID,
					},
					After: &models.ProductVariant{
						Name:                 rO[iOO],
						ProductVariantTypeID: pvts[iOV].ID,
					},
				})
			}
		}

		//* product variant medias - check changes
		// validation
		// if iteration 0
		if iOV == 0 {
			// if length of images not equal options then invalid model
			if len(rV.Images) != len(rV.Options) {
				return nil, e.ErrInvalidVariantsMedia
			}
		} else {
			// if iteration more than 0 then don't allow image to exists
			if len(rV.Images) > 0 {
				return nil, e.ErrInvalidVariantsMedia
			}
		}
		// make array of uint of product variant ids
		pvIDs := make([]uint, 0)
		for _, v := range pv {
			pvIDs = append(pvIDs, v.ID)
		}
		// get all product variant medias from db by
		pvm, err := s.productVariantMediaRepository.GetAllByProductVariantID(pvIDs)
		if err != nil && err != e.ErrDataNotFound {
			return nil, err
		}
		// product variant medias - req to ori comparison - append creates
		if len(rV.Images) > len(oV.Images) {
			for i := len(oV.Images); i < len(rV.Images); i++ {
				pvmCreates = append(pvmCreates, &productInternalDtos.PvmUpdateInfo{
					After: &models.ProductVariantMedia{
						MediaURL: rV.Images[i],
					},
				})
			}
		}
		// product variant medias - ori to req comparison - append updates and deletes
		// ? how to get variant id?
		for iOI, oI := range oV.Images {
			rI := rV.Images
			// if index not exist then append deletes
			if len(rI)-1 < iOI {
				pvmDeletes = append(pvmDeletes, &productInternalDtos.PvmUpdateInfo{
					Before: &models.ProductVariantMedia{
						ID: pvm[iOI].ID,
					},
				})
				continue
			} else if oI != rI[iOI] {
				// if there is changes then append updates
				pvmUpdates = append(pvmUpdates, &productInternalDtos.PvmUpdateInfo{
					Before: &models.ProductVariantMedia{
						ID: pvm[iOI].ID,
					},
					After: &models.ProductVariantMedia{
						MediaURL: rI[iOI],
					},
				})
			}
		}
	}

	// validation - product items length
	if len(request.ProductItems) != piLengthValidation {
		return nil, e.ErrInvalidModel
	}

	var piCreates = []*productInternalDtos.PiUpdateInfo{}
	var piUpdates = []*productInternalDtos.PiUpdateInfo{}
	var piDeletes = []*productInternalDtos.PiUpdateInfo{}

	rPI := request.ProductItems
	oPI := oriProductModel.ProductItems

	//* product items - check changes
	// product items - req to ori comparison - append creates
	if len(rPI) > len(oPI) {
		for i := len(oPI); i < len(rPI); i++ {
			piCreates = append(piCreates, &productInternalDtos.PiUpdateInfo{
				After: &models.ProductItem{
					ProductID:     product.ID,
					SKU:           rPI[i].SKU,
					StockQuantity: rPI[i].StockQuantity,
					Price:         rPI[i].Price,
				},
			})
		}
	}
	// product itmes - ori to req comparison - append updates and deletes
	for iOPI, vOPI := range oPI {
		// if index not exist then append deletes
		if len(rPI)-1 < iOPI {
			piDeletes = append(piDeletes, &productInternalDtos.PiUpdateInfo{
				Before: &models.ProductItem{
					ID: vOPI.ID,
				},
			})
			continue
		}

		// if there are no changes
		if vOPI.Price == rPI[iOPI].Price &&
			vOPI.StockQuantity == rPI[iOPI].StockQuantity &&
			vOPI.SKU == rPI[iOPI].SKU {
			continue
		}

		piUpdates = append(piUpdates, &productInternalDtos.PiUpdateInfo{
			Before: &models.ProductItem{
				ID: vOPI.ID,
			},
			After: &models.ProductItem{
				Price:         rPI[iOPI].Price,
				StockQuantity: rPI[iOPI].StockQuantity,
				SKU:           rPI[iOPI].SKU,
			},
		})
	}

	var pcsCreates = []*productInternalDtos.PcsUpdateInfo{}
	var pcsDeletes = []*productInternalDtos.PcsUpdateInfo{}

	//* product courier services - check changes
	// original couriers
	oC := oriProductModel.Couriers
	rC := request.Couriers
	// validate couriers length
	if len(oC) != len(rC) {
		return nil, e.ErrInvalidModel
	}
	for iVOC, vOC := range oC {
		oCS := vOC.CourierServices
		rCS := request.Couriers[iVOC].CourierServices
		// validate courier services length
		if len(oCS) != len(rCS) {
			return nil, e.ErrInvalidModel
		}
		for iVOCS, vOCS := range oCS {
			merchantCourierService, err := s.merchantCourierServiceRepository.GetByMerchantIDCourierServiceID(oriProductModel.MerchantID, vOCS.ID)
			if err != nil {
				return nil, e.ErrInvalidModel
			}

			// if nothing changes then continue to next iteration
			if vOCS.IsSelected == rCS[iVOCS].IsSelected {
				continue
			}
			// if is active change from true to false then delete product courier services
			if vOCS.IsSelected && !rCS[iVOCS].IsSelected {
				// append to delete
				pcsDeletes = append(pcsDeletes, &productInternalDtos.PcsUpdateInfo{
					Before: &models.ProductCourierService{
						MerchantCourierServiceID: merchantCourierService.ID,
						ProductId:                product.ID,
					},
				})
			}
			if !vOCS.IsSelected && rCS[iVOCS].IsSelected {
				// append to create
				pcsCreates = append(pcsCreates, &productInternalDtos.PcsUpdateInfo{
					After: &models.ProductCourierService{
						MerchantCourierServiceID: merchantCourierService.ID,
						ProductId:                product.ID,
					},
				})
			}
		}
	}

	//* product medias - update operations
	for _, v := range pmCreates {
		fmt.Println("pmCreates")
		fmt.Println("\tproduct_id:", v.After.ProductID)
		fmt.Println("\tmedia_url:", v.After.MediaURL)
		_, err := s.productMediaRepository.Create(v.After)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range pmDeletes {
		fmt.Println("pmDeletes")
		fmt.Println("\tproduct_id:", v.Before.ProductID)
		fmt.Println("\tmedia_url:", v.Before.MediaURL)
		fmt.Println(v)
		err := s.productMediaRepository.DeleteByProductIDMediaURLType(v.Before.ProductID, v.Before.MediaURL, v.Before.Type)
		if err != nil {
			return nil, err
		}
	}

	//* product variant types - update operations
	for _, v := range pvtCreates {
		fmt.Println("pvtCreates", v.NewPVT.Name)
		_, err := s.productVariantTypeRepository.Create(v.NewPVT)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range pvtUpdates {
		fmt.Println("pvtUpdates", v.ID)
		_, err := s.productVariantTypeRepository.Update(v.ID, v.NewPVT)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range pvtDeletes {
		fmt.Println("pvtDeletes", v.ID)
		_, err := s.productVariantTypeRepository.Delete(v.ID)
		if err != nil {
			return nil, err
		}
	}

	//* product variant - update operations
	newPvIds := make([]uint, 0)
	for _, v := range pvCreates {
		fmt.Println("pvtCreates", v.After.Name)
		newPv, err := s.productVariantRepository.Create(v.After)
		if err != nil {
			return nil, err
		}
		newPvIds = append(newPvIds, newPv.ID)
	}
	for _, v := range pvUpdates {
		fmt.Println("pvUpdates", v.Before.Name, v.After.Name)
		_, err := s.productVariantRepository.Update(v.Before.ID, v.After)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range pvDeletes {
		fmt.Println("pvDeletes", v.Before.ID)
		_, err := s.productVariantRepository.Delete(v.Before.ID)
		if err != nil {
			return nil, err
		}
	}

	//* product variant medias - update operations
	for i, v := range pvmCreates {
		fmt.Println("pvmCreates", v.After.MediaURL)
		v.After.ProductVariantID = newPvIds[i]
		_, err := s.productVariantMediaRepository.Create(v.After)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range pvmUpdates {
		fmt.Println("pvmUpdates", v.Before.ID, v.After.MediaURL)
		_, err := s.productVariantMediaRepository.Update(v.Before.ID, v.After)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range pvmDeletes {
		// since pvm dependence on pv id then when pv deletes above then it wil cascade delete too to its pvm, thats why we add condition if err is not not found then its not an error
		_, err = s.productVariantMediaRepository.Delete(v.Before.ID)
		if err != nil && err != e.ErrDataNotFound {
			return nil, err
		}
	}

	//* product items - update operations
	for _, v := range piCreates {
		fmt.Println("piCreates")
		fmt.Printf("\tprice: %v\n\tstock quantity: %v\n\tsku: %v\n", v.After.Price, v.After.StockQuantity, v.After.SKU)
		_, err := s.productItemRepository.CreateProductItem(v.After)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range piUpdates {
		fmt.Println("piUpdates")
		fmt.Printf("\tid: %v\n\tprice: %v\n\tstock quantity: %v\n\tsku: %v\n", v.Before.ID, v.After.Price, v.After.StockQuantity, v.After.SKU)
		newPiMap := map[string]any{
			"price":          v.After.Price,
			"stock_quantity": v.After.StockQuantity,
			"sku":            v.After.SKU,
		}
		// to allow updates stock and SKU to 0 use maps for updates
		err := s.productItemRepository.UpdateUsingMap(v.Before.ID, newPiMap)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range piDeletes {
		fmt.Println("piDeletes")
		fmt.Println("\tid:", v.Before.ID)
		err = s.productItemRepository.Delete(v.Before.ID)
		if err != nil && err != e.ErrDataNotFound {
			return nil, err
		}
	}

	//* product courier services - update creates
	for _, v := range pcsCreates {
		fmt.Println("pcsCreates")
		fmt.Println("\t", v.After.MerchantCourierServiceID)
		fmt.Println("\t", v.After.ProductId)
		_, err = s.productCourierServiceRepository.Create(v.After)
		if err != nil && err != e.ErrDataNotFound {
			return nil, err
		}

	}
	for _, v := range pcsDeletes {
		fmt.Println("pcsDeletes")
		fmt.Println("\t", v.Before.MerchantCourierServiceID)
		fmt.Println("\t", v.Before.ProductId)
		err = s.productCourierServiceRepository.DeleteByMerchantCourierServiceIDProductID(v.Before.MerchantCourierServiceID, v.Before.ProductId)
		if err != nil && err != e.ErrDataNotFound {
			return nil, err
		}
	}

	// var pVarCreates = []*productInternalDtos.PvarUpdateInfo{}
	//* product variations - map creates
	productItems, err := s.productItemRepository.GetAllByProductID(product.ID)
	if err != nil {
		return nil, err
	}
	// get all product variants as index
	productVariantsType, err := s.productVariantRepository.GetAllVariantIndexByProductID(product.ID)
	if err != nil {
		return nil, err
	}
	variantsMap := make([][]uint, len(productVariantsType))
	for i, v := range productVariantsType {
		variantsMap[i] = make([]uint, len(v.Variants))
		for iV, vV := range v.Variants {
			variantsMap[i][iV] = uint(vV)
		}
	}
	var pVarCreates = make([]*productInternalDtos.PvarUpdateInfo, 0)
	// variant indexes combination map for validate and making sure combination not duplicate
	var viCombinationMap = map[string]bool{}
	for iPI, vPI := range request.ProductItems {
		// validate variant index length to request variant length (variant type length)
		if len(vPI.VariantIndex) != len(request.Variants) {
			return nil, e.ErrInvalidModel
		}
		choosenVariantIndexes := make([]uint, 0)
		// variant index combination map to append later
		var viCombination string
		for iVI, vVI := range vPI.VariantIndex {
			// validate if value is less than request variant length (variant type length)
			// if request variant index is more than variant type length then err
			if int(vVI) > len(variantsMap) {
				return nil, e.ErrInvalidModel
			}
			// append choosen variant indexes
			choosenVariantIndexes = append(choosenVariantIndexes, variantsMap[iVI][vVI])
			// append combination
			viCombination += strconv.Itoa(int(vVI))
		}
		// validate if combination already added before
		if exists := viCombinationMap[viCombination]; exists {
			return nil, e.ErrInvalidModel
		}
		// add combination to map
		viCombinationMap[viCombination] = true

		// append create
		for _, vCVI := range choosenVariantIndexes {
			pVarCreates = append(pVarCreates, &productInternalDtos.PvarUpdateInfo{
				After: &models.ProductVariation{
					ProductItemID:    productItems[iPI].ID,
					ProductVariantID: vCVI,
				},
			})
		}
	}

	//* delete all product variation
	for _, oPI := range oriProductModel.ProductItems {
		err := s.productVariationRepository.DeleteAllByProductItemID(oPI.ID)
		if err != nil && err != e.ErrDataNotFound {
			return nil, err
		}
	}
	for _, v := range pVarCreates {
		fmt.Println("pVarCreates")
		fmt.Println("\t", v.After.ProductItemID)
		fmt.Println("\t", v.After.ProductVariantID)
		_, err := s.productVariationRepository.Create(v.After)
		if err != nil {
			return nil, err
		}
	}

	// get updated product model
	productModel, err := s.GetUpdateModelByProductID(product.ID)
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
	for iC, vC := range couriers {
		courierServices, err := s.merchantCourierServiceRepository.GetAllCourierServicesForProductUpdate(vC.ID, productModel.MerchantID, productModel.ID)
		if err != nil && err != e.ErrDataNotFound {
			return nil, err
		}

		csIsSelected := false
		for _, cs := range courierServices {
			if cs.IsSelected {
				csIsSelected = true
			}
		}

		productModel.Couriers[iC].CourierServices = courierServices
		productModel.Couriers[iC].IsSelected = csIsSelected
	}

	return productModel, nil
}
func (s *productService) GetAvailableProductsByPromotionPeriod(id uuid.UUID, query *dtos.AvailableProductsByPromotionPeriodQuery) ([]*dtos.AvailableProductsResponse, error) {
	merchant, err := s.merchantRepository.GetByUserID(id)

	if err != nil {
		return nil, err
	}

	availableProduct, err := s.productRepository.GetAvailableProductsByPromotionPeriod(merchant.ID, query)

	if err != nil {
		return nil, err
	}

	return dtos.NewAvailableProductsResponse(availableProduct), nil
}

func (s *productService) CheckCourierService(request *productInternalDtos.ProductCheckCourierService) (*productInternalDtos.ProductCheckCourierServiceResponse, error) {
	return s.productRepository.CheckCourierService(request)
}

func (s *productService) MyShopGetAll(query *productInternalDtos.ProductQuery) (responseProductDetails []*dtos.ProductDetailsGetByIDResponse, pagination *utils.Pagination, err error) {
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

	products, pagination, err := s.productRepository.GetAll(query, pagination)
	if err != nil {
		return
	}

	for _, vP := range products {
		query := &internalDtos.ProductDetailsInternalQuery{
			ProductID: vP.ID,
			Status:    vP.ProductStatus,
		}
		productDetails, err := s.ProductDetailsGetByID(query)
		if err != nil {
			return nil, nil, err
		}
		responseProductDetails = append(responseProductDetails, productDetails)
	}

	return responseProductDetails, pagination, nil
}

func (s *productService) GetAllUserFavoriteProducts(userID uuid.UUID, query *productInternalDtos.ProductFavoriteQuery) (productDetails []*productInternalDtos.Product, pagination *utils.Pagination, err error) {
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

	productDetails, pagination, err = s.productRepository.GetAllUserFavoriteProducts(userID, query, pagination)
	if err != nil {
		return
	}

	return productDetails, pagination, nil
}

func (s *productService) DeleteFavoriteProducts(userID uuid.UUID, request *dtos.DeleteFavoriteProductsRequest) ([]*dtos.UpdateProductFavoriteResponse, error) {
	var productFavorite []*models.ProductFavorite

	for _, productID := range request.ProductIDs {
		favoriteProduct, err := s.productFavoriteRepository.GetByUserIDAndProductID(userID, *productID)

		if err != nil {
			return nil, err
		}

		favoriteProduct.IsFavorite = false

		updatedFavoriteProduct, err := s.productFavoriteRepository.Update(favoriteProduct.ID, favoriteProduct)

		if err != nil {
			return nil, err
		}

		productFavorite = append(productFavorite, updatedFavoriteProduct)
	}

	return dtos.NewBatchUpdateProductFavoriteResponse(productFavorite), nil
}
