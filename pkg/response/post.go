package response

import (
	"backend/pkg/models"
	"time"
)

type Post struct {
	// A post is a message in a ring.
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`

	// Ring attributes
	RingName  string `json:"ringName"`
	RingColor string `json:"ringColor"`

	AuthorUsername string       `json:"authorUsername" gorm:"index"`
	Author         *models.User `json:"author"`
	Title          string       `json:"title"`
	Body           string       `json:"body,omitempty"`
	Link           *string      `json:"link"`
	Domain         *string      `json:"domain" gorm:"index"`
	Score          int          `json:"score"`
	CommentsCount  int          `json:"commentsCount"`
	Ups            int          `json:"ups"`
	Downs          int          `json:"downs"`
	Nsfw           bool         `json:"nsfw"`
	VotedUp        bool         `json:"votedUp"`
	VotedDown      bool         `json:"votedDown"`
}
