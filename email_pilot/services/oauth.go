package services

import (
	"context"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GetOAuthConfig returns the Google OAuth configuration.
func GetOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GMAIL_CLIENT_ID"),
		ClientSecret: os.Getenv("GMAIL_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.readonly"},
		Endpoint:     google.Endpoint,
	}
}

// GetClient returns an HTTP client using the given token.
func GetClient(ctx context.Context, config *oauth2.Config, token *oauth2.Token) *http.Client {
	return config.Client(ctx, token)
}
