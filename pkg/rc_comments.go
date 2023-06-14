package server

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func (s *Server) getRcComments(c *gin.Context) {
	id := c.Param("id")
	parts := strings.SplitN(id, ".", 2)
	if len(parts) != 2 || parts[1] != "json" {
		badRequest(c)
		return
	}
	postId, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		s.logger.Warnf("invalid post id: %s, %v", id, err)
		badRequest(c)
		return
	}

	// Get Post:
	post, err := s.repoPost(uint(postId))
	if err != nil {
		internalServerError(c)
		return
	}

	comments, err := s.repoComments(uint(postId))
	if err != nil {
		internalServerError(c)
		return
	}

	redditComments, err := toRedditComments(post, comments)
	if err != nil {
		internalServerError(c)
		return
	}
	c.JSON(200, redditComments)
}
