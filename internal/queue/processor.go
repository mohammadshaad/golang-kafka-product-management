package queue

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/mohammadshaad/zocket/pkg/util"
)

var s3Client *util.S3Client

func InitS3Storage(bucketName string) {
	client, err := util.InitS3Client(bucketName)
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}
	s3Client = client
}

func ProcessImageMessage(key, value []byte) error {
	url := string(value)
	log.Printf("Processing image for URL: %s", url)

	// Step 1: Download Image
	img, err := util.DownloadImage(url)
	if err != nil {
		return fmt.Errorf("error downloading image: %w", err)
	}

	// Step 2: Compress Image
	compressed, err := util.CompressImage(img, 75)
	if err != nil {
		return fmt.Errorf("error compressing image: %w", err)
	}

	// Step 3: Upload Compressed Image to S3
	fileName := filepath.Base(url) // Use the original filename
	filePath := fmt.Sprintf("compressed/%s", fileName)
	s3URL, err := s3Client.UploadToS3(filePath, compressed)
	if err != nil {
		return fmt.Errorf("error uploading to S3: %w", err)
	}

	log.Printf("Successfully processed and uploaded image: %s", s3URL)
	return nil
}
