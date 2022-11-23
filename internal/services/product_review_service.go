package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
)

type ProductReviewService interface {
	Create(*dtos.ProductReviewCreateRequest) (*dtos.ProductReviewCreateResponse, error)
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

func (pr *productReviewService) Create(request *dtos.ProductReviewCreateRequest) (*dtos.ProductReviewCreateResponse, error) {
	productReview := &models.ProductReview{
		UserID:        request.UserID,
		OrderDetailID: request.OrderDetailID,
		ProductID:     request.ProductID,
		Rating:        request.Rating,
		Review:        request.Review,
	}
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
		}
		pr.productReviewMediaRepository.Create(productReviewMedia)
	}
	response := &dtos.ProductReviewCreateResponse{
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
