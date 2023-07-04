package server

import (
	"backend/pkg/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func (s *Server) getPost(c *gin.Context) {
	id, done := parsePostId(c)
	if done {
		return
	}

	post, err := s.repoPost(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{"error": "post not found"})
		return
	}

	if err != nil {
		s.logger.Errorf("unable to get post %d: %v", id, err)
		internalServerError(c)
		return
	}

	// If user is logged, check if user has upvoted this post
	username := c.GetString("username")
	if username != "" {
		action, err := s.repoPostAction(uint(id), username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// User hasn't voted this post
				c.JSON(200, post)
				return
			} else {
				s.logger.Errorf("unable to check if user %s has upvoted post %d: %v", username, id, err)
				internalServerError(c)
				return
			}
		}
		post.VotedUp = action.Action == models.ActionUpvote
		post.VotedDown = action.Action == models.ActionDownvote
	}

	c.JSON(200, post)
}

func parsePostId(c *gin.Context) (int64, bool) {
	return parseId(c, "id")
}
func parseId(c *gin.Context, name string) (int64, bool) {
	idParam := c.Param(name)
	if idParam == "" {
		c.JSON(400, gin.H{"error": name + " is required"})
		return 0, true
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": name + " must be a number"})
		return 0, true
	}

	if id < 0 {
		c.JSON(400, gin.H{"error": name + " must be a positive number"})
		return 0, true
	}
	return id, false
}
