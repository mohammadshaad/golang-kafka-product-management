package util

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

// UploadToS3 uploads a file to the S3 bucket
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

	url := "https://" + s.bucketName + ".s3.amazonaws.com/" + fileName
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
		// Compress JPEG image
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		if err != nil {
			log.Printf("Error compressing JPEG image: %v", err)
			return nil, err
		}
	case color.RGBAModel, color.NRGBAModel:
		// PNG compression (no quality param needed for PNG)
		err := png.Encode(&buf, img)
		if err != nil {
			log.Printf("Error compressing PNG image: %v", err)
			return nil, err
		}
	default:
		// If it's neither a JPEG nor PNG, return an error
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
