package server

import (
	"backend/pkg/models"
	"backend/pkg/request"
	"fmt"
	"net/url"
)

func (s *Server) repoCreatePost(req request.PostCreate, username string) (*models.Post, error) {
	// Check if ring exists
	ringExists, err := s.repoRingExists(req.Ring)
	if err != nil {
		return nil, fmt.Errorf("unable to check if ring exists: %v", err)
	}

	if !ringExists {
		return nil, fmt.Errorf(ErrRingDoesNotExist)
	}

	// Validate Post Request
	err = req.Validate()
	if err != nil {
		s.logger.Warnf("invalid post create request: %v", err)
		return nil, fmt.Errorf(ErrInvalidPostRequest)
	}

	// Create post
	post := models.Post{
		Title:          req.Title,
		Nsfw:           req.Nsfw,
		RingName:       req.Ring,
		AuthorUsername: username,
	}

	if req.Link != nil {
		post.Link = req.Link

		// Determine domain
		domain, err := getDomain(*req.Link)
		if err != nil {
			return nil, fmt.Errorf("unable to get domain: %v", err)
		}
		post.Domain = &domain
	}

	if req.Body != nil {
		post.Body = *req.Body
	}

	err = s.db.Create(&post).Error
	if err != nil {
		return nil, fmt.Errorf("unable to create post: %v", err)
	}
	return &post, nil
}

func getDomain(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("invalid link")
	}
	return u.Hostname(), nil
}
