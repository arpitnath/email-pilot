package models

import "time"

// OAuthToken represents the structure for storing OAuth tokens.
type OAuthToken struct {
	ID           uint      `gorm:"primaryKey"`
	UserEmail    string    `gorm:"unique;not null"`
	AccessToken  string    `gorm:"not null"`
	RefreshToken string    `gorm:"not null"`
	Expiry       time.Time `gorm:"not null"`
}
