package models

import "time"

type Ring struct {
	// A ring is a community: a collection of posts.
	Name         string    `json:"name" gorm:"primaryKey"`
	Title        string    `json:"title"`
	DisplayName  string    `json:"displayName"`
	Description  string    `json:"description"`
	Posts        []Post    `json:"posts,omitempty" gorm:"foreignKey:RingName;references:Name"`
	CreatedOn    time.Time `json:"createdOn" gorm:"autoCreateTime"`
	Nsfw         bool      `json:"nsfw"`
	Subscribers  int       `json:"subscribers"`
	PrimaryColor string    `json:"primaryColor"`
}
