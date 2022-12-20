package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
)

type ImageService interface {
	Upload(uploadImageRequest *externalDtos.UploadImageRequest) (uploadImageResponse *externalDtos.UploadImageResponse, err error)
	Delete(deleteImageRequest *externalDtos.DeleteImageRequest) (deleteImageResponse *externalDtos.DeleteImageResponse, err error)
}

type imageService struct {
	// empty struct
}

type ImageServiceConfig struct {
	// empty struct
}

func NewImageService(c *ImageServiceConfig) ImageService {
	return &imageService{}
}

func (i *imageService) Upload(uploadImageRequest *externalDtos.UploadImageRequest) (uploadImageResponse *externalDtos.UploadImageResponse, err error) {
	photo, err := uploadImageRequest.Photo.Open()
	if err != nil {
		return
	}
	imageUrl, err := utils.UploadImage(photo)
	uploadImageResponse = &externalDtos.UploadImageResponse{
		Photo: imageUrl,
	}
	return
}

func (i *imageService) Delete(deleteImageRequest *externalDtos.DeleteImageRequest) (deleteImageResponse *externalDtos.DeleteImageResponse, err error) {
	result, err := utils.DeleteImage(deleteImageRequest.PublicID)
	if err != nil {
		return
	}
	deleteImageResponse = &externalDtos.DeleteImageResponse{
		Result: result,
	}
	return
}
