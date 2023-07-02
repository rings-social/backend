package request

import (
	"errors"
	"regexp"
)

type CreateRingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

var colorRegexp = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func (r CreateRingRequest) Validate() error {
	if len(r.Title) > 100 {
		return errors.New("title too long")
	}

	if len(r.Description) > 1000 {
		return errors.New("description too long")
	}

	if !colorRegexp.MatchString(r.Color) {
		return errors.New("invalid color")
	}

	return nil
}
