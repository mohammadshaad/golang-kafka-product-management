package util

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

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
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		log.Printf("Error compressing image: %v", err)
		return nil, err
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
