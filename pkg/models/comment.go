package models

type Comment struct {
	// A comment is a reply to a post.
	Model
	Post     *Post    `json:"post,omitempty" gorm:"foreignKey:PostId"`
	PostId   uint     `json:"post_id"`
	Parent   *Comment `json:"parent,omitempty" gorm:"foreignKey:ParentId"`
	ParentId *uint    `json:"parent_id"`

	// The author of the comment.
	AuthorUsername string `json:"author_id"`
	Author         User   `json:"author" gorm:"foreignKey:AuthorUsername;references:Username"`

	// The comment's content.
	Body  string `json:"body"`
	Ups   uint   `json:"ups" gorm:"default:0"`
	Downs uint   `json:"downs"`
	Score int    `json:"score"`

	Depth   int       `json:"depth" gorm:"-"`
	Replies []Comment `json:"replies" gorm:"-"`
}
