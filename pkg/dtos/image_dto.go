package dtos

import "mime/multipart"

type UploadImageRequest struct {
	Photo multipart.FileHeader `form:"photo" binding:"required"`
}

type UploadImageResponse struct {
	Photo string `json:"photo"`
}
