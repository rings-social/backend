package models

import "time"

type PostAction struct {
	Username  string `json:"username" gorm:"primaryKey"`
	User      User   `json:"user,omitempty" gorm:"foreignKey:Username"`
	Post      Post   `json:"post,omitempty"`
	PostId    uint   `json:"post_id" gorm:"primaryKey"`
	Action    string `json:"action" gorm:"type:post_action;index"`
	CreatedAt time.Time
}
