// @title Savannah Store API
// @version 1.0
// @description This is the API documentation for Savannah Store.
// @host localhost:8080
// @BasePath /api

package main

import (
	"log"

	docs "github.com/Oj-washingtone/savannah-store/docs"
	"github.com/Oj-washingtone/savannah-store/internal/api"
	"github.com/Oj-washingtone/savannah-store/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	database.ConnectDB()
	router := gin.Default()

	apiGroup := router.Group("/api")
	api.AppRoutes(apiGroup)

	// Serve Swagger UI
	docs.SwaggerInfo.BasePath = "/api"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":8080")
}
