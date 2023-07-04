package models

type Post struct {
	// A post is a message in a ring.
	Model
	RingName       string  `json:"ringName" gorm:"index"`
	Ring           *Ring   `json:"ring,omitempty" gorm:"foreignKey:RingName;references:Name"`
	AuthorUsername string  `json:"author_username" gorm:"index"`
	Author         *User   `json:"author,omitempty" gorm:"foreignKey:AuthorUsername;references:Username"`
	Title          string  `json:"title"`
	Body           string  `json:"body,omitempty"`
	Link           *string `json:"link"`
	Domain         *string `json:"domain" gorm:"index"`
	Score          int     `json:"score"`
	CommentsCount  int     `json:"commentsCount"`
	Ups            int     `json:"ups"`
	Downs          int     `json:"downs"`
	Nsfw           bool    `json:"nsfw"`

	VotedUp   bool `json:"votedUp" gorm:"-"`
	VotedDown bool `json:"votedDown" gorm:"-"`
}
