package queue

import (
    "encoding/json"
    "fmt"
    "log"
    "path/filepath"

    "github.com/mohammadshaad/zocket/pkg/util"
)

type ImageMessage struct {
    ProductID int    `json:"product_id"`
    ImageURL  string `json:"image_url"`
}

var s3Client *util.S3Client

func InitS3Storage(bucketName string) {
    client, err := util.InitS3Client(bucketName)
    if err != nil {
        log.Fatalf("Failed to initialize S3 client: %v", err)
    }
    s3Client = client
}

func ProcessImageMessage(key, value []byte) error {
    // Parse the message
    var msg ImageMessage
    if err := json.Unmarshal(value, &msg); err != nil {
        return fmt.Errorf("error unmarshaling message: %w", err)
    }

    log.Printf("Processing image for product %d: %s", msg.ProductID, msg.ImageURL)

    // Download Image
    img, err := util.DownloadImage(msg.ImageURL)
    if err != nil {
        return fmt.Errorf("error downloading image: %w", err)
    }

    // Compress Image
    compressed, err := util.CompressImage(img, 75)
    if err != nil {
        return fmt.Errorf("error compressing image: %w", err)
    }

    // Generate unique filename for the compressed image
    originalFilename := filepath.Base(msg.ImageURL)
    compressedFilename := fmt.Sprintf("compressed/%d_%s", msg.ProductID, originalFilename)

    // Upload Compressed Image to S3
    s3URL, err := s3Client.UploadToS3(compressedFilename, compressed)
    if err != nil {
        return fmt.Errorf("error uploading to S3: %w", err)
    }

    // Update the product record with the compressed image URL
    if err := util.UpdateProductImageURL(msg.ProductID, msg.ImageURL, s3URL); err != nil {
        return fmt.Errorf("error updating product image URL: %w", err)
    }

    log.Printf("Successfully processed image for product %d. Compressed URL: %s", msg.ProductID, s3URL)
    return nil
}
