package models

import "time"

type Post struct {
	// A post is a message in a ring.
	ID             int       `json:"id" gorm:"primaryKey"`
	PostedOn       time.Time `json:"posted_on"`
	RingName       string    `json:"ring_name" gorm:"index"`
	Ring           *Ring     `json:"ring,omitempty" gorm:"foreignKey:RingName;references:Name"`
	AuthorUsername string    `json:"author_username" gorm:"index"`
	Author         *User     `json:"author,omitempty" gorm:"foreignKey:AuthorUsername;references:Username"`
	Title          string    `json:"title"`
	Body           string    `json:"body,omitempty"`
	Link           string    `json:"link"`
	Domain         string    `json:"domain" gorm:"index"`
	Score          int       `json:"score"`
}
