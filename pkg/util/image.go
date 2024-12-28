package util

import (
    "bytes"
    "context"
    "image"
    "image/color"
    "image/jpeg"
    "image/png"
    "fmt"
	"strconv"
    "log"
    "mime"
    "net/http"
    "os"
    "path/filepath"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"

    "github.com/mohammadshaad/zocket/internal/db"
	"github.com/mohammadshaad/zocket/internal/cache"
)

// S3Client wraps the AWS S3 client
type S3Client struct {
    client     *s3.Client
    bucketName string
}

// InitS3Client initializes the S3 client with the specified bucket name
func InitS3Client(bucketName string) (*S3Client, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Printf("Error loading AWS configuration: %v", err)
        return nil, err
    }

    client := s3.NewFromConfig(cfg)
    return &S3Client{client: client, bucketName: bucketName}, nil
}

// UploadToS3 uploads a file to the S3 bucket and returns the URL
func (s *S3Client) UploadToS3(fileName string, fileData []byte) (string, error) {
    contentType := mime.TypeByExtension(filepath.Ext(fileName))
    if contentType == "" {
        contentType = "application/octet-stream"
    }

    input := &s3.PutObjectInput{
        Bucket:      aws.String(s.bucketName),
        Key:         aws.String(fileName),
        Body:        bytes.NewReader(fileData),
        ContentType: aws.String(contentType),
    }

    _, err := s.client.PutObject(context.TODO(), input)
    if err != nil {
        log.Printf("Error uploading to S3: %v", err)
        return "", err
    }

    url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, fileName)
    log.Printf("File uploaded to S3: %s", url)
    return url, nil
}

// DownloadImage downloads an image from a given URL
func DownloadImage(url string) (image.Image, error) {
    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Error downloading image: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    img, _, err := image.Decode(resp.Body)
    if err != nil {
        log.Printf("Error decoding image: %v", err)
        return nil, err
    }

    return img, nil
}

// CompressImage compresses an image and returns a byte array
func CompressImage(img image.Image, quality int) ([]byte, error) {
    var buf bytes.Buffer
    switch img.ColorModel() {
    case color.YCbCrModel:
        err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
        if err != nil {
            log.Printf("Error compressing JPEG image: %v", err)
            return nil, err
        }
    case color.RGBAModel, color.NRGBAModel:
        err := png.Encode(&buf, img)
        if err != nil {
            log.Printf("Error compressing PNG image: %v", err)
            return nil, err
        }
    default:
        log.Printf("Unsupported image format: %T", img)
        return nil, fmt.Errorf("unsupported image format: %T", img)
    }

    return buf.Bytes(), nil
}

// SaveImageToFile saves a compressed image to the local file system
func SaveImageToFile(imgData []byte, filePath string) error {
    file, err := os.Create(filePath)
    if err != nil {
        log.Printf("Error creating file: %v", err)
        return err
    }
    defer file.Close()

    _, err = file.Write(imgData)
    if err != nil {
        log.Printf("Error writing file: %v", err)
        return err
    }

    return nil
}

// UpdateProductImageURL updates the product record with the new compressed image URL
func UpdateProductImageURL(productID int, originalURL, compressedURL string) error {
    if db.DB == nil {
        return fmt.Errorf("database connection not initialized")
    }

    var product db.Product
    if err := db.DB.First(&product, productID).Error; err != nil {
        log.Printf("Error finding product %d: %v", productID, err)
        return err
    }

    // Find the original image URL in the ProductImages array and get its index
    originalIndex := -1
    for i, url := range product.ProductImages {
        if url == originalURL {
            originalIndex = i
            break
        }
    }

    if originalIndex == -1 {
        return fmt.Errorf("original image URL not found in product images")
    }

    // Ensure CompressedProductImages array is initialized and has the same length
    for len(product.CompressedProductImages) < len(product.ProductImages) {
        product.CompressedProductImages = append(product.CompressedProductImages, "")
    }

    // Update the compressed image URL at the corresponding index
    product.CompressedProductImages[originalIndex] = compressedURL

    // Save the updated product
    if err := db.DB.Save(&product).Error; err != nil {
        log.Printf("Error saving product with compressed image URL: %v", err)
        return err
    }

    // Invalidate the cache for this product
    if err := cache.InvalidateProductCache(strconv.Itoa(productID)); err != nil {
        log.Printf("Error invalidating cache for product %d: %v", productID, err)
    }

    log.Printf("Successfully updated product %d with compressed image URL", productID)
    return nil
}
