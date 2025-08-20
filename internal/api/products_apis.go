package api

import (
	"github.com/Oj-washingtone/savannah-store/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")

	{
		// product categories

		products.POST("/categories/create", handlers.AddProductCategory)
		products.GET("/categories", handlers.ListCategories)
		products.PATCH("/categories/:id", handlers.UpdateCategory)

		// products
		products.POST("/create", handlers.CreateProduct)
		products.GET("/:id", handlers.GetProductById)
		products.GET("/", handlers.ListProducts)
		products.PATCH("/:id", handlers.UpdateProduct)
		products.DELETE("/:id", handlers.DeleteProduct)

	}
}
