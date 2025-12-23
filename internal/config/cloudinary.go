package config

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var Cld *cloudinary.Cloudinary

func InitCloudinary() {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		log.Println("Warning: CLOUDINARY_URL not set, image upload will not be available")
		return
	}

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
	Cld = cld
	log.Println("Cloudinary initialized successfully")
}

func UploadImage(ctx context.Context, file interface{}, publicID string) (string, error) {
	if Cld == nil {
		return "", errors.New("cloudinary not initialized")
	}
	
	resp, err := Cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: publicID,
		Folder:   "ewsbe/news",
	})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}
