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

// maybeAuthenticatedUser is like authenticatedUser, but it doesn't
// require the user to be authenticated. If the user is authenticated, it
// sets the "username" context variable.
func (s *Server) maybeAuthenticatedUser(c *gin.Context) {
	if !s.hasIdToken(c) {
		return
	}

	idToken, done := s.idToken(c)
	if done {
		return
	}

	username, err := s.usernameForIdToken(idToken)
	if err != nil {
		return
	}

	c.Set("username", username)
}
