package services

import (
	"fmt"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/constants/product_media_constants"
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
	ProductDetailsGetByID(productID uint) (*dtos.ProductDetailsGetByIDResponse, error)
	GetPurchaseQuantitiesByID(userID uuid.UUID, quantity *dtos.PurchaseQuantitiesQuery) (*dtos.GetPurchaseQuantitiesByIDResponse, error)
	ViewProduct(userID uuid.UUID, productID uint) (*dtos.GetProductViewResponse, error)
	GetFavoriteProduct(userID uuid.UUID, productID uint) (*dtos.GetProductFavoriteResponse, error)
	FavoriteProduct(userID uuid.UUID, productID uint) (*dtos.GetProductFavoriteResponse, error)
	CreateProduct(request *productInternalDtos.ProductModel) (*productInternalDtos.ProductModel, error)
	UpdateProductUsingModel(userID uuid.UUID, request *productInternalDtos.ProductModel) (*productInternalDtos.ProductModel, error)
	GetUpdateModelByProductID(productID uint) (*productInternalDtos.ProductModel, error)
}

type productService struct {
	productRepository                r.ProductRepository
	merchantRepository               r.MerchantRepository
	productItemRepository            r.ProductItemRepository
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
	return &u
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
	// productVariants := []*models.ProductVariantType{}
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
	fmt.Println(oriProductModel.Description)

	// product validaiton
	// check if product status not its enum
	if request.Status != constants.ProductStatusActive &&
		request.Status != constants.ProductStatusInactive {
		return nil, e.ErrInvalidModel
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
	var pmCreate = []string{}
	// product medias to delete
	var pmDelete = []string{}

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
			pmCreate = append(pmCreate, v)
		}
	}
	// add product media to delete
	for _, v := range oriProductModel.Images {
		if exists := reqIPMMap[v]; !exists {
			pmDelete = append(pmDelete, v)
		}
	}

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
			pmCreate = append(pmCreate, v)
		}
	}
	// add product media to delete
	for _, v := range oriProductModel.Videos {
		if exists := reqVPMMap[v]; !exists {
			pmDelete = append(pmDelete, v)
		}
	}

	fmt.Println("pmCreate", pmCreate)
	fmt.Println("pmDelete", pmDelete)

	// todo: create product medias
	// todo: delete product medias

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
	if err != nil {
		return nil, err
	}

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
						Name:                 oO,
						ProductVariantTypeID: pvts[iOV].ID,
					},
				})
				continue
			} else if oO != rO[iOO] {
				// if there is changes then append updates
				pvUpdates = append(pvUpdates, &productInternalDtos.PvUpdateInfo{
					Before: &models.ProductVariant{
						Name:                 oO,
						ProductVariantTypeID: pvts[iOV].ID,
					},
					After: &models.ProductVariant{
						Name:                 rO[iOO],
						ProductVariantTypeID: pvts[iOV].ID,
					},
				})
			}
		}

		// todos: check images (product variants media) changes
		// product variant medias - req to ori comparison - append creates
		if len(rV.Images) > len(oV.Images) {
			for i := len(oV.Images); i < len(rV.Images); i++ {
				pvmCreates = append(pvmCreates, &productInternalDtos.PvmUpdateInfo{
					After: &models.ProductVariantMedia{
						MediaURL: rV.Images[i],
						// todo: fill product variant id later after inserting product variant id
					},
				})
			}
		}
		// product ariant medias - ori to req comparison - append updates and deletes
		// ? how to get variant id?
		for iOI, oI := range oV.Images {
			rI := rV.Images
			// if index not exist then append deletes
			if len(rI)-1 < iOI {
				pvmDeletes = append(pvmDeletes, &productInternalDtos.PvmUpdateInfo{
					Before: &models.ProductVariantMedia{
						MediaURL: oI,
					},
				})
				continue
			} else if oI != rI[iOI] {
				// if there is changes then append updates
				pvmUpdates = append(pvmUpdates, &productInternalDtos.PvmUpdateInfo{
					Before: &models.ProductVariantMedia{
						MediaURL: oI,
					},
					After: &models.ProductVariantMedia{
						MediaURL: rI[iOI],
					},
				})
			}
		}
	}

	var piCreates = []*productInternalDtos.PiUpdateInfo{}
	var piUpdates = []*productInternalDtos.PiUpdateInfo{}
	var piDeletes = []*productInternalDtos.PiUpdateInfo{}

	rPI := request.ProductItems
	oPI := oriProductModel.ProductItems

	// todo: check changes on product item
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
			pvmDeletes = append(pvmDeletes, &productInternalDtos.PvmUpdateInfo{
				Before: &models.ProductVariantMedia{
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
		// ! to allow updates stock and SKU to 0 quantity then don't forget to use maps for updates later
	}

	// todo: check changes on product variation

	// todo: check changes on product courier service

	// todo: get created product model

	// todo: get updated product model

	// log product variant types update info
	for iPvt, pvt := range pvtCreates {
		fmt.Println("pvtCreates", iPvt, pvt.NewPVT.Name)
		// delete pvt also whatever that referencing to pvt id
	}
	for iPvtU, pvtU := range pvtUpdates {
		fmt.Println("pvtUpdates", iPvtU, pvtU.ID, pvtU.NewPVT.Name)
		// update pvt
	}
	for iPvtD, pvtD := range pvtDeletes {
		fmt.Println("pvtDeletes", iPvtD, pvtD.ID)
		// delete pvt also whatever that referencing to pvt id
	}

	// log product variants update info
	for _, v := range pvCreates {
		fmt.Println("pvCreates", v.After.Name)
	}
	for _, v := range pvDeletes {
		fmt.Println("pvDeletes", v.Before.Name)
	}
	for _, v := range pvUpdates {
		fmt.Println("pvUpdates", v.Before.Name, v.After.Name)
	}

	// log product variant medias update info
	for _, v := range pvmCreates {
		fmt.Println("pvmCreates", v.After.MediaURL)
	}
	for _, v := range pvmDeletes {
		fmt.Println("pvmDeletes", v.Before.MediaURL)
	}
	for _, v := range pvmUpdates {
		fmt.Println("pvmUpdates", v.Before.MediaURL, v.After.MediaURL)
	}

	// log product items update info
	for _, v := range piCreates {
		fmt.Println("piCreates")
		fmt.Printf("\tid: %v\n\tprice: %v\n\tstock quantity: %v\n\tsku: %v\n", v.Before.ID, v.After.Price, v.After.StockQuantity, v.After.SKU)
	}
	for _, v := range piDeletes {
		fmt.Println("> piDeletes")
		fmt.Println(v.After.Price)
		fmt.Println(v.After.StockQuantity)
		fmt.Println(v.After.SKU)
	}
	for _, v := range piUpdates {
		fmt.Println("piUpdates")
		fmt.Printf("\tid: %v\n\tprice: %v\n\tstock quantity: %v\n\tsku: %v\n", v.Before.ID, v.After.Price, v.After.StockQuantity, v.After.SKU)
	}

	productModel, err := s.GetUpdateModelByProductID(product.ID)
	if err != nil {
		return nil, err
	}

	return productModel, nil
}

func (s *productService) GetUpdateModelByProductID(productID uint) (*productInternalDtos.ProductModel, error) {
	// TODO: compare product model

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
