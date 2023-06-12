package models

import "time"

type Ring struct {
	// A ring is a community: a collection of posts.
	Name        string    `json:"name" gorm:"primaryKey"`
	Description string    `json:"description"`
	Posts       []Post    `json:"posts,omitempty" gorm:"foreignKey:RingName;references:Name"`
	CreatedOn   time.Time `json:"created_on" gorm:"autoCreateTime"`
}
