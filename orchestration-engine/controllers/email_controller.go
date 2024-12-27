package controllers

import (
	"email_pilot/config"
	"email_pilot/models"
	"email_pilot/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func GetEmailsController(c *gin.Context) {
	var token models.OAuthToken
	if err := config.DB.First(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch token"})
		return
	}

	emailService := services.EmailService{
		Token: &oauth2.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			Expiry:       token.Expiry,
		},
		OauthConfig: services.GetOAuthConfig(),
	}

	emails, err := emailService.FetchEmails(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch emails: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"emails": emails})
}
