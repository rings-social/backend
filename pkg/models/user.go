package models

import (
	"gorm.io/gorm"
	"time"
)

type SocialLink struct {
	// A social link is a link to a user's profile on another site.
	Username string `json:"username" gorm:"primaryKey"`
	Platform string `json:"platform" gorm:"primaryKey"`
	Url      string `json:"url"`
}

type User struct {
	// A user is a person who can post to rings.
	Username    string         `json:"username" gorm:"primaryKey"`
	DisplayName string         `json:"displayName"`
	SocialLinks []SocialLink   `json:"socialLinks,omitempty" gorm:"foreignKey:Username;references:Username"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
	Posts       []Post         `json:"posts,omitempty" gorm:"foreignKey:AuthorUsername;references:Username"`
}
