package main

import (
	"github.com/Oj-washingtone/savannah-store/internal/api"
	"github.com/Oj-washingtone/savannah-store/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {

	// connect to database
	database.ConnectDB()

	// start application
	router := gin.Default()
	api.AppRoutes(router)
	router.Run()
}
