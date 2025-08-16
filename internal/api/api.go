package api

import (
	"github.com/Oj-washingtone/savannah-store/internal/handlers"
	"github.com/gin-gonic/gin"
)

func AppRoutes(router *gin.Engine) {
	router.GET("/albums", handlers.GetAlbums)
	router.POST("/albums", handlers.AddToAlbum)
	router.GET("/albums/:id", handlers.GetAlbumById)
}
