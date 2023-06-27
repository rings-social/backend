package models

import "time"

type VoteAction = string

const ActionUpvote VoteAction = "upvote"
const ActionDownvote VoteAction = "downvote"

type CommentAction struct {
	Username  string  `json:"username" gorm:"primaryKey"`
	User      User    `json:"user,omitempty" gorm:"foreignKey:Username"`
	Comment   Comment `json:"comment,omitempty"`
	CommentId uint    `json:"comment_id" gorm:"primaryKey"`
	Action    string  `json:"action" gorm:"type:comment_action;index"`
	CreatedAt time.Time
}
