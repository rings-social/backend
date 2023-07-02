package server

import (
	"backend/pkg/request"
	"github.com/gin-gonic/gin"
)

func (s *Server) createPost(c *gin.Context) {
	var postCreateReq request.PostCreate
	err := c.BindJSON(&postCreateReq)
	if err != nil {
		s.logger.Warnf("unable to bind json: %v", err)
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// Get current user
	username := c.GetString("username")
	post, err := s.repoCreatePost(postCreateReq, username)

	if err != nil {
		if err.Error() == ErrRingDoesNotExist {
			c.JSON(400, gin.H{"error": "ring does not exist"})
			return
		}

		if err.Error() == ErrInvalidPostRequest {
			c.JSON(400, gin.H{"error": "invalid post request"})
			return
		}

		s.logger.Errorf("unable to create post: %v", err)
		internalServerError(c)
		return
	}

	c.JSON(200, post)
}
