package server

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

/*
comments: ID of the comment tree to return
*/
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

	parentIdParam := c.Query("comment")
	var parentId *uint
	if parentIdParam != "" {
		parentIdInt, err := strconv.ParseInt(parentIdParam, 10, 64)
		if err != nil {
			s.logger.Warnf("invalid parent id: %s, %v", parentIdParam, err)
			badRequest(c)
			return
		}
		parentIdUint := uint(parentIdInt)
		parentId = &parentIdUint
	}

	comments, err := s.repoComments(uint(postId), parentId)
	if err != nil {
		internalServerError(c)
		return
	}

	for k, comment := range comments {
		childComments, err := s.repoComments(uint(postId), &comment.ID)
		if err != nil {
			s.logger.Errorf("unable to get child comments for comment %d: %v", comment.ID, err)
			internalServerError(c)
			return
		}
		comments[k].Replies = childComments
	}

	redditComments, err := toRedditComments(post, comments)
	if err != nil {
		internalServerError(c)
		return
	}
	c.JSON(200, redditComments)
}
