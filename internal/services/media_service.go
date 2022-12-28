package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

type MediaService interface {
	Upload(uploadMediaRequest *externalDtos.UploadMediaRequest) (uploadMediaResponse *externalDtos.UploadMediaResponse, err error)
	Delete(deleteMediaRequest *externalDtos.DeleteMediaRequest) (deleteMediaResponse *externalDtos.DeleteMediaResponse, err error)
}

type mediaService struct {
	// empty struct
}

type MediaServiceConfig struct {
	// empty struct
}

func NewMediaService(c *MediaServiceConfig) MediaService {
	return &mediaService{}
}

func (i *mediaService) Upload(uploadMediaRequest *externalDtos.UploadMediaRequest) (uploadMediaResponse *externalDtos.UploadMediaResponse, err error) {
	mediaSize := uploadMediaRequest.Media.Size / 1000000
	mediaContentType := uploadMediaRequest.Media.Header.Values("Content-Type")[0]
	if mediaContentType == "image/jpeg" && mediaSize > 5 {
		err = errors.ErrFileSize
		return
	} else if mediaContentType == "video/mp4" && mediaSize > 10 {
		err = errors.ErrFileSize
		return
	}
	media, err := uploadMediaRequest.Media.Open()
	if err != nil {
		return
	}
	defer media.Close()

	mediaUrl, err := utils.UploadMedia(media)
	uploadMediaResponse = &externalDtos.UploadMediaResponse{
		Media: mediaUrl,
	}
	return
}

func (i *mediaService) Delete(deleteMediaRequest *externalDtos.DeleteMediaRequest) (deleteMediaResponse *externalDtos.DeleteMediaResponse, err error) {
	result, err := utils.DeleteMedia(deleteMediaRequest.PublicID)
	if err != nil {
		return
	}
	deleteMediaResponse = &externalDtos.DeleteMediaResponse{
		Result: result,
	}
	return
}
