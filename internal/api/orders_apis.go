package api

import (
	"github.com/Oj-washingtone/savannah-store/internal/handlers"
	"github.com/Oj-washingtone/savannah-store/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterOrdersRoutes(router *gin.RouterGroup) {
	orders := router.Group("/orders")

	{
		orders.POST("/create", middleware.AuthMiddleware(), handlers.CreateOrder)
		orders.GET("/", handlers.GetAllOrders)
	}
}
