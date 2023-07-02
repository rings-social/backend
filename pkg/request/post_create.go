package request

import (
	"fmt"
	"net/url"
)

type PostCreate struct {
	Title string  `json:"title"`
	Body  *string `json:"body,omitempty"`
	Link  *string `json:"link,omitempty"`
	Nsfw  bool    `json:"nsfw"`
	Ring  string  `json:"ring"`
}

func (c PostCreate) Validate() error {
	// Titles can be maximum 300 characters
	if len(c.Title) > 300 {
		return fmt.Errorf("title cannot be longer than 60 characters")
	}

	// Body can be maximum 1000 characters
	if c.Body != nil && len(*c.Body) > 1000 {
		return fmt.Errorf("body cannot be longer than 1000 characters")
	}

	// Link can be maximum 300 characters
	if c.Link != nil && len(*c.Link) > 300 {
		return fmt.Errorf("link cannot be longer than 300 characters")
	}

	// Parse link
	if c.Link != nil {
		u, err := url.Parse(*c.Link)
		if err != nil {
			return fmt.Errorf("invalid link")
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return fmt.Errorf("invalid link scheme")
		}
	}

	if c.Link != nil && c.Body != nil {
		return fmt.Errorf("cannot have both link and body")
	}

	return nil
}
