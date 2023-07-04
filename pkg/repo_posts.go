package server

import (
	"backend/pkg/models"
	"backend/pkg/request"
	"errors"
	"fmt"
	"gorm.io/gorm"
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

func (s *Server) repoVoteAction(action models.VoteAction, username string, id int64) error {
	// Check if post exists
	post, err := s.repoPost(uint(id))
	if err != nil {
		s.logger.Warnf("unable to get post: %v", err)
		return fmt.Errorf(ErrUnableToGetPost)
	}

	// Check if user has already voted
	vote, err := s.repoGetVote(username, uint(id))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Warnf("unable to get vote: %v", err)
		return fmt.Errorf(ErrUnableToGetVote)
	}

	if vote != nil {
		if vote.Action == action {
			// User has already voted with this action
			return fmt.Errorf(ErrUnableToVoteUserAlreadyVoted)
		}

		vote.Action = action

		// Run in a transaction
		tx := s.db.Begin()
		err = tx.Save(&vote).Error
		if err != nil {
			s.logger.Warnf("unable to save vote: %v", err)
			tx.Rollback()
			return fmt.Errorf(ErrUnableToSaveVote)
		}

		increase := 2
		if action == models.ActionDownvote {
			increase = -2
		}
		// Vote saved, update post, add 2 to score (1 for upvote, 1 for removing downvote)
		err = s.repoIncreasePostScore(tx, post.ID, increase)
		if err != nil {
			s.logger.Warnf("unable to increase post score by %d: %v", increase, err)
			tx.Rollback()
			return fmt.Errorf(ErrUnableToIncreasePostScore)
		}
		tx.Commit()
	} else {
		// User has not voted yet
		vote := models.PostAction{
			Username: username,
			PostId:   uint(id),
			Action:   action,
		}

		// Run in a transaction
		tx := s.db.Begin()
		err = tx.Create(&vote).Error
		if err != nil {
			s.logger.Warnf("unable to create vote: %v", err)
			tx.Rollback()
			return fmt.Errorf(ErrUnableToCreateVote)
		}
		// Vote saved, update post, add +-1 to score
		increase := 1
		if action == models.ActionDownvote {
			increase = -1
		}
		err = s.repoIncreasePostScore(tx, post.ID, increase)
		if err != nil {
			s.logger.Warnf("unable to increase post score by %d: %v", increase, err)
			tx.Rollback()
			return fmt.Errorf(ErrUnableToIncreasePostScore)
		}
		tx.Commit()
	}

	return nil
}
