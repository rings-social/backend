package models

import (
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
	Username       string       `json:"username" gorm:"primaryKey"`
	DisplayName    string       `json:"displayName"`
	ProfilePicture *string      `json:"profilePicture"`
	SocialLinks    []SocialLink `json:"socialLinks,omitempty" gorm:"foreignKey:Username;references:Username"`
	CreatedAt      time.Time    `json:"createdAt" gorm:"autoCreateTime"`
	DeletedAt      *time.Time   `json:"deletedAt,omitempty"`
	Posts          []Post       `json:"posts,omitempty" gorm:"foreignKey:AuthorUsername;references:Username"`

	// A user might have multiple badges
	Badges []Badge `json:"badges,omitempty" gorm:"many2many:user_badges;"`

	AuthSubject *string `json:"-" gorm:"uniqueIndex"`
	Admin       bool    `json:"admin"`
}
