package testutils

import (
    "github.com/gin-gonic/gin"
    "github.com/mohammadshaad/zocket/internal/api"
    "github.com/mohammadshaad/zocket/internal/db"
)

// TestProduct is a mock product for testing
var TestProduct = db.Product{
    UserID:             1,
    ProductName:        "Test Product",
    ProductDescription: "Test Description",
    ProductImages:      []string{"https://shaad-my-product-images-bucket.s3.eu-north-1.amazonaws.com/shaad-image.jpeg"},
    ProductPrice:       69.69,
}

// SetupTestRouter initializes a test router
func SetupTestRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    api.SetupRoutes(router)
    return router
}