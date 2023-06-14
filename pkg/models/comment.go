package models

import "gorm.io/gorm"

type Comment struct {
	// A comment is a reply to a post.
	gorm.Model
	Post     Post     `json:"post" gorm:"foreignKey:PostId"`
	PostId   uint     `json:"postId"`
	Parent   *Comment `json:"parent" gorm:"foreignKey:ParentId"`
	ParentId *uint    `json:"parentId"`

	// The author of the comment.
	AuthorUsername string `json:"authorId"`
	Author         User   `json:"author" gorm:"foreignKey:AuthorUsername;references:Username"`

	// The comment's content.
	Body  string `json:"body"`
	Ups   uint   `json:"ups" gorm:"default:0"`
	Downs uint   `json:"downs"`
}
