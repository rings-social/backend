package models

import "time"

type SocialLink struct {
	// A social link is a link to a user's profile on another site.
	Username string `json:"username" gorm:"primaryKey"`
	Platform string `json:"platform" gorm:"primaryKey"`
	Url      string `json:"url"`
}

type User struct {
	// A user is a person who can post to rings.
	Username    string       `json:"username" gorm:"primaryKey"`
	DisplayName string       `json:"display_name"`
	SocialLinks []SocialLink `json:"social_links,omitempty" gorm:"foreignKey:Username;references:Username"`
	CreatedOn   time.Time    `json:"created_on" gorm:"autoCreateTime"`
	Posts       []Post       `json:"posts,omitempty" gorm:"foreignKey:AuthorUsername;references:Username"`
}
