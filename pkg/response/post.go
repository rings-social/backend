package response

import "time"

type Post struct {
	// A post is a message in a ring.
	ID       int       `json:"id" gorm:"primaryKey"`
	PostedOn time.Time `json:"posted_on"`

	// Ring attributes
	RingName  string `json:"ring_name"`
	RingColor string `json:"ring_color"`

	AuthorUsername string  `json:"author_username" gorm:"index"`
	Title          string  `json:"title"`
	Body           string  `json:"body,omitempty"`
	Link           *string `json:"link"`
	Domain         *string `json:"domain" gorm:"index"`
	Score          int     `json:"score"`
	CommentsCount  int     `json:"commentsCount"`
	Ups            int     `json:"ups"`
	Downs          int     `json:"downs"`
	Nsfw           bool    `json:"nsfw"`
}
