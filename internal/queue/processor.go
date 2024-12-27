package queue

import (
	"fmt"
	"log"
	"path/filepath"
	"github.com/mohammadshaad/zocket/pkg/util"
)

func ProcessImageMessage(key, value []byte) error {
	url := string(value)
	log.Printf("Processing image for URL: %s", url)

	// Step 1: Download Image
	img, err := util.DownloadImage(url)
	if err != nil {
		return err
	}

	// Step 2: Compress Image
	compressed, err := util.CompressImage(img, 75)
	if err != nil {
		return err
	}

	// Step 3: Save Compressed Image
	outputPath := filepath.Join("output", fmt.Sprintf("%s.jpg", key))
	err = util.SaveImageToFile(compressed, outputPath)
	if err != nil {
		return err
	}

	log.Printf("Successfully processed and saved image: %s", outputPath)
	return nil
}
