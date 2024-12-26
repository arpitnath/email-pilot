package controllers

import (
	"context"
	"email_pilot/config"
	"email_pilot/models"
	"email_pilot/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// AuthController handles the OAuth authentication flow.
func AuthController(c *gin.Context) {
	oauthConfig := services.GetOAuthConfig()

	// Step 1: Redirect the user to Google's OAuth 2.0 consent screen.
	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// CallbackController handles the callback from Google's OAuth 2.0.
func CallbackController(c *gin.Context) {
	oauthConfig := services.GetOAuthConfig()

	// Step 2: Exchange the authorization code for an access token.
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is missing"})
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Failed to exchange authorization code for token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Step 3: Save the token in the database for later use.
	oauthToken := &models.OAuthToken{
		UserEmail:    "arpitnath@gmail.com",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
	if err := config.DB.Create(oauthToken).Error; err != nil {
		log.Printf("Failed to save token to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token saved successfully"})
}
