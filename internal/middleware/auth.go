package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/Oj-washingtone/savannah-store/internal/authenticator"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		rawIDToken := parts[1]

		auth, err := authenticator.New()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to init authenticator"})
			c.Abort()
			return
		}

		oidcConfig := &oidc.Config{
			ClientID: os.Getenv("AUTH0_CLIENT_ID"),
		}

		idToken, err := auth.Verifier(oidcConfig).Verify(c.Request.Context(), rawIDToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse token claims"})
			c.Abort()
			return
		}

		c.Set("user", claims)

		c.Next()
	}
}
