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
