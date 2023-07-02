package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) authenticatedUser(c *gin.Context) {
	if !s.hasIdToken(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	idToken, done := s.idToken(c)
	if done {
		return
	}

	// Make sure the user exists in the database
	username, err := s.usernameForIdToken(idToken)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	c.Set("username", username)
}
