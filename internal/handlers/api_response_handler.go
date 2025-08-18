package handlers

import "github.com/gin-gonic/gin"

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func RespondSuccess(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func RespondError(c *gin.Context, status int, message string, err string) {
	c.JSON(status, ApiResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}
