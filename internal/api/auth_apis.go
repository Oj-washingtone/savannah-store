package api

import (
	"net/http"

	"github.com/Oj-washingtone/savannah-store/internal/authenticator"
	"github.com/Oj-washingtone/savannah-store/internal/handlers"
	"github.com/Oj-washingtone/savannah-store/internal/model"
	"github.com/Oj-washingtone/savannah-store/internal/repocitory"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var claims struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func RegisterAuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")

	{
		auth.GET("/login", handlers.Login)

		auth.GET("/auth0/callback", func(c *gin.Context) {
			code := c.Query("code")

			if code == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "no code in query"})
				return
			}

			auth, err := authenticator.New()

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to init authenticator"})
				return
			}

			token, err := auth.Exchange(c.Request.Context(), code)

			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to exchange code for token"})
				return
			}

			idToken, err := auth.VerifyIDToken(c.Request.Context(), token)

			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid id token"})
				return
			}

			if err := idToken.Claims(&claims); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse claims"})
				return
			}

			userRepo := repocitory.NewUserRepository()

			user, err := userRepo.GetByEmail(c.Request.Context(), claims.Email)

			if err != nil {
				newUser := &model.User{
					Name:    claims.Name,
					Email:   claims.Email,
					Auth0Id: claims.Sub,
					Role:    "customer",
				}

				newUser.ID = uuid.New()

				if err := userRepo.Create(c.Request.Context(), newUser); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "failed to create user",
						"details": err.Error(),
					})
					return
				}

				user = newUser

			}

			c.IndentedJSON(http.StatusOK, gin.H{
				"user":  user,
				"token": token.Extra("id_token"),
			})
		})
	}

}
