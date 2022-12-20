package dtos

import "mime/multipart"

type UploadImageRequest struct {
	Photo multipart.FileHeader `form:"photo" binding:"required"`
}

type UploadImageResponse struct {
	Photo string `json:"photo"`
}

type DeleteImageRequest struct {
	PublicID string `json:"public_id" binding:"required"`
}

type DeleteImageResponse struct {
	Result string `json:"result"`
}
