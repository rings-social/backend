package models

import (
	"time"
)

type Ring struct {
	// A ring is a community: a collection of posts.
	Name          string     `json:"name" gorm:"primaryKey"`
	Title         string     `json:"title"`
	DisplayName   string     `json:"displayName"`
	Description   string     `json:"description"`
	Posts         []Post     `json:"posts,omitempty" gorm:"foreignKey:RingName;references:Name"`
	CreatedAt     time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
	Nsfw          bool       `json:"nsfw"`
	PrimaryColor  string     `json:"primaryColor"`
	OwnerUsername string     `json:"ownerUsername"`
	Owner         *User      `json:"owner,omitempty" gorm:"foreignKey:OwnerUsername;references:Username"`
	Subscribers   uint64     `json:"subscribers"`
}
