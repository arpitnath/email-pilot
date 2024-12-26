package controllers

import (
	"email_pilot/config"
	"email_pilot/models"
	"email_pilot/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// GetEmailsController handles the /emails endpoint
func GetEmailsController(c *gin.Context) {
	// Retrieve the token from the database
	var token models.OAuthToken
	if err := config.DB.First(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch token"})
		return
	}

	// Initialize EmailService
	emailService := services.EmailService{
		Token: &oauth2.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			Expiry:       token.Expiry,
		},
		OauthConfig: services.GetOAuthConfig(),
	}

	// Fetch emails
	emails, err := emailService.FetchEmails(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch emails: " + err.Error()})
		return
	}

	// Return emails in response
	c.JSON(http.StatusOK, gin.H{"emails": emails})
}
