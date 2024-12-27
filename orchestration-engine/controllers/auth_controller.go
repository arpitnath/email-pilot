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

func AuthController(c *gin.Context) {
	oauthConfig := services.GetOAuthConfig()

	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func CallbackController(c *gin.Context) {
	oauthConfig := services.GetOAuthConfig()

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

	oauthToken := &models.OAuthToken{
		UserEmail:    "arpitnath@gmail.com", // TODO: get user email from google: just for testing
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
