package api

import (
	"github.com/Oj-washingtone/savannah-store/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")

	{
		products.POST("/categories", handlers.AddProductCategory)
		products.GET("/categories", handlers.GetProductCategories)
	}
}
