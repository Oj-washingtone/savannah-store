package api

import (
	"github.com/Oj-washingtone/savannah-store/internal/handlers"
	"github.com/Oj-washingtone/savannah-store/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(router *gin.RouterGroup) {
	cart := router.Group("/cart")

	cart.Use(middleware.AuthMiddleware())

	{
		cart.POST("/create", handlers.AddToCart)
		cart.DELETE("/remove/:id", handlers.RemoveFromCart)
		cart.GET("/", handlers.GetCartItems)
		cart.PATCH("/update/quantity/:id", handlers.UpdateQuantity)
	}
}
