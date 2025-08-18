package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/Oj-washingtone/savannah-store/internal/authenticator"
	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary Login endpoint
// @Description Redirects the user to the Auth0/Google login page for authentication.
// @Tags auth
// @Produce json
// @Success 302 {string} string "Redirects to Auth0/Google login page"
// @Failure 500 {object} map[string]string "Failed to initialize authenticator"
// @Router /auth/login [get]
func Login(c *gin.Context) {

	auth, err := authenticator.New()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to initialize authenticator", "error": err.Error()})
		return
	}

	state, err := generateRandomState()

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	redirectURL := auth.AuthCodeURL(state)

	c.Redirect(http.StatusFound, redirectURL)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
