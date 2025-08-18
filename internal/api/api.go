package api

import (
	"github.com/gin-gonic/gin"
)

func AppRoutes(router *gin.RouterGroup) {
	RegisterAuthRoutes(router)
}
