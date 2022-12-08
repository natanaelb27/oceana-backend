package services

import (
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
)

type ProductReviewService interface {
	Create(*externalDtos.ProductReviewCreateRequest) (*externalDtos.ProductReviewCreateResponse, error)
	GetByProductID(id uint, query *externalDtos.ProductReviewQuery) (*externalDtos.ProductReviewGetByProductIDResponse, *utils.Pagination, error)
}

type productReviewService struct {
	productReviewRepository      r.ProductReviewRepository
	productReviewMediaRepository r.ProductReviewMediaRepository
}

type ProductReviewServiceConfig struct {
	ProductReviewRepository      r.ProductReviewRepository
	ProductReviewMediaRepository r.ProductReviewMediaRepository
}

func NewProductReviewService(c *ProductReviewServiceConfig) ProductReviewService {
	return &productReviewService{
		productReviewRepository:      c.ProductReviewRepository,
		productReviewMediaRepository: c.ProductReviewMediaRepository,
	}
}

func (pr *productReviewService) Create(request *externalDtos.ProductReviewCreateRequest) (*externalDtos.ProductReviewCreateResponse, error) {
	productReview := request.ToModel()
	newProductReview, err := pr.productReviewRepository.Create(productReview)
	if err != nil {
		return nil, err
	}
	if request.Photo1 != nil {
		photoFile, err := request.Photo1.Open()
		if err != nil {
			return nil, err
		}
		imageUrl, err := utils.UploadImage(photoFile)
		if err != nil {
			return nil, err
		}
		productReviewMedia := &models.ProductReviewMedia{
			ProductReviewID: newProductReview.ID,
			MediaURL:        imageUrl,
			Type:            "photo",
		}
		pr.productReviewMediaRepository.Create(productReviewMedia)
	}
	if request.Photo2 != nil {
		photoFile, err := request.Photo2.Open()
		if err != nil {
			return nil, err
		}
		imageUrl, err := utils.UploadImage(photoFile)
		if err != nil {
			return nil, err
		}
		productReviewMedia := &models.ProductReviewMedia{
			ProductReviewID: newProductReview.ID,
			MediaURL:        imageUrl,
			Type:            "photo",
		}
		pr.productReviewMediaRepository.Create(productReviewMedia)
	}
	if request.Photo3 != nil {
		photoFile, err := request.Photo3.Open()
		if err != nil {
			return nil, err
		}
		imageUrl, err := utils.UploadImage(photoFile)
		if err != nil {
			return nil, err
		}
		productReviewMedia := &models.ProductReviewMedia{
			ProductReviewID: newProductReview.ID,
			MediaURL:        imageUrl,
			Type:            "photo",
		}
		pr.productReviewMediaRepository.Create(productReviewMedia)
	}
	if request.Photo4 != nil {
		photoFile, err := request.Photo4.Open()
		if err != nil {
			return nil, err
		}
		imageUrl, err := utils.UploadImage(photoFile)
		if err != nil {
			return nil, err
		}
		productReviewMedia := &models.ProductReviewMedia{
			ProductReviewID: newProductReview.ID,
			MediaURL:        imageUrl,
			Type:            "photo",
		}
		pr.productReviewMediaRepository.Create(productReviewMedia)
	}
	if request.Photo5 != nil {
		photoFile, err := request.Photo5.Open()
		if err != nil {
			return nil, err
		}
		imageUrl, err := utils.UploadImage(photoFile)
		if err != nil {
			return nil, err
		}
		productReviewMedia := &models.ProductReviewMedia{
			ProductReviewID: newProductReview.ID,
			MediaURL:        imageUrl,
			Type:            "photo",
		}
		pr.productReviewMediaRepository.Create(productReviewMedia)
	}
	if request.Video != nil {
		videoFile, err := request.Video.Open()
		if err != nil {
			return nil, err
		}
		videoUrl, err := utils.UploadImage(videoFile)
		if err != nil {
			return nil, err
		}
		productReviewMedia := &models.ProductReviewMedia{
			ProductReviewID: newProductReview.ID,
			MediaURL:        videoUrl,
			Type:            "video",
		}
		pr.productReviewMediaRepository.Create(productReviewMedia)
	}
	response := &externalDtos.ProductReviewCreateResponse{
		ID:            newProductReview.ID,
		UserID:        newProductReview.UserID,
		OrderDetailID: newProductReview.OrderDetailID,
		ProductID:     newProductReview.ProductID,
		Rating:        newProductReview.Rating,
		Review:        newProductReview.Review,
		CreatedAt:     newProductReview.CreatedAt,
		UpdatedAt:     newProductReview.UpdatedAt,
	}
	return response, nil
}

func (pr *productReviewService) GetByProductID(id uint, query *externalDtos.ProductReviewQuery) (*externalDtos.ProductReviewGetByProductIDResponse, *utils.Pagination, error) {
	pagination := &utils.Pagination{
		Limit: query.Limit,
		Page:  query.Page,
	}
	productReviews, pagination, err := pr.productReviewRepository.GetByProductID(id, pagination)
	if err != nil {
		return nil, nil, err
	}
	var productReviewDataList []*internalDtos.ProductReviewGetByProductID
	for _, productReview := range productReviews {
		var productReviewMediaLinks []*string
		productReviewMedia, err := pr.productReviewMediaRepository.GetByProductReviewID(productReview.ID)
		if err != nil {
			return nil, nil, err
		}
		for _, singleProductReviewMedia := range productReviewMedia {
			productReviewMediaLinks = append(productReviewMediaLinks, &singleProductReviewMedia.MediaURL)
		}
		productReviewData := internalDtos.NewProductReviewGetByProductID(productReview, productReviewMediaLinks)
		productReviewDataList = append(productReviewDataList, productReviewData)
	}
	response := externalDtos.NewProductReviewGetByIDResponse(productReviewDataList)
	return response, pagination, nil
}
