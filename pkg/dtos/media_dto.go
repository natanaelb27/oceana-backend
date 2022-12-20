package dtos

import "mime/multipart"

type UploadMediaRequest struct {
	Media multipart.FileHeader `form:"media" binding:"required"`
}

type UploadMediaResponse struct {
	Media string `json:"media_url"`
}

type DeleteMediaRequest struct {
	PublicID string `json:"public_id" binding:"required"`
}

type DeleteMediaResponse struct {
	Result string `json:"result"`
}
