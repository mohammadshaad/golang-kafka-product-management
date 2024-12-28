package api

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/mohammadshaad/zocket/internal/db"
    "github.com/mohammadshaad/zocket/internal/queue"
)

func AddProductHandler(c *gin.Context) {
    var product db.Product

    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Initialize CompressedProductImages array
    product.CompressedProductImages = make([]string, len(product.ProductImages))

    if err := db.DB.Create(&product).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save product"})
        return
    }

    // Publish each image to Kafka for processing
    for _, url := range product.ProductImages {
        msg := queue.ImageMessage{
            ProductID: int(product.ID),
            ImageURL:  url,
        }
        
        msgBytes, err := json.Marshal(msg)
        if err != nil {
            log.Printf("Error marshaling image message: %v", err)
            continue
        }

        if err := queue.PublishMessage([]byte(strconv.Itoa(int(product.ID))), msgBytes); err != nil {
            log.Printf("Failed to enqueue image: %v", err)
            continue
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Product added successfully",
        "product": product,
    })
}

func GetProductByIDHandler(c *gin.Context) {
    id := c.Param("id")

    var product db.Product
    if err := db.DB.First(&product, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, product)
}

func GetAllProductsHandler(c *gin.Context) {
    userID := c.Query("user_id")
    minPrice := c.Query("min_price")
    maxPrice := c.Query("max_price")

    var products []db.Product
    query := db.DB

    if userID != "" {
        query = query.Where("user_id = ?", userID)
    }

    if minPrice != "" {
        min, _ := strconv.ParseFloat(minPrice, 64)
        query = query.Where("product_price >= ?", min)
    }

    if maxPrice != "" {
        max, _ := strconv.ParseFloat(maxPrice, 64)
        query = query.Where("product_price <= ?", max)
    }

    if err := query.Find(&products).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
        return
    }

    c.JSON(http.StatusOK, products)
}
