package server

import (
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

	c.JSON(200, post)
}

func parsePostId(c *gin.Context) (int, bool) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return 0, true
	}

	// Parse ID as uint
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "id must be a number"})
		return 0, true
	}

	if id < 0 {
		c.JSON(400, gin.H{"error": "id must be a positive number"})
		return 0, true
	}
	return id, false
}

func (s *Server) getComments(c *gin.Context) {
	id, done := parsePostId(c)
	if done {
		return
	}

	var parentId *uint
	parentIdParam := c.Query("parent_id")
	if parentIdParam != "" {
		parentIdInt, err := strconv.Atoi(parentIdParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "parent_id must be a number"})
			return
		}
		if parentIdInt < 0 {
			c.JSON(400, gin.H{"error": "parent_id must be a positive number"})
			return
		}

		parentIdUint := uint(parentIdInt)
		parentId = &parentIdUint
	}

	comments, err := s.repoComments(uint(id), parentId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{"error": "comments not found"})
		return
	}

	if err != nil {
		s.logger.Errorf("unable to get comments for post %d: %v", id, err)
		internalServerError(c)
		return
	}

	c.JSON(200, comments)
}
