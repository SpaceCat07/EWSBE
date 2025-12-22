package config

import (
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var Cld *cloudinary.Cloudinary

func InitCloudinary() {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		log.Fatal("CLOUDINARY_URL is required")
	}

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
	Cld = cld
}

func UploadImage(ctx context.Context, file interface{}, publicID string) (string, error) {
	resp, err := Cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: publicID,
		Folder:   "ewsbe/news",
	})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}
