package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohammadshaad/zocket/internal/db"
	"github.com/mohammadshaad/zocket/internal/queue"
)

func AddProductHandler (c *gin.Context) {
	var product db.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save product"})
		return
	}

	for _, url := range product.ProductImages {
		if err := queue.PublishMessage(nil, []byte(url)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue image"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added successfully"})
}


func GetProductByIDHandler (c *gin.Context) {
	id := c.Param("id")

	var product db.Product

	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func GetAllProductsHandler (c *gin.Context) {
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