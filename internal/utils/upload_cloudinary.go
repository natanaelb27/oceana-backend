package utils

import (
	"context"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImage(input interface{}) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()

	cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))
	useFilenameValue := true
	resp, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{Folder: os.Getenv("CLOUDINARY_UPLOAD_FOLDER"), UseFilename: &useFilenameValue})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}

func DeleteImage(publicID string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()
	cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))
	resp, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	if err != nil {
		return "", err
	}
	return resp.Result, nil
}
