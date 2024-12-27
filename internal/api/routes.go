package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes (router *gin.Engine) {
	api := router.Group("/api/v1")

	{
		api.POST("/products", AddProductHandler)
		api.GET("/products/:id", GetProductByIDHandler)
		api.GET("/products", GetAllProductsHandler)
	}
}