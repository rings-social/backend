package server

import (
	"backend/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) innerRouteVotePost(c *gin.Context, action models.VoteAction) {
	postId, done := parsePostId(c)
	if done {
		return
	}

	username := c.GetString("username")
	err := s.repoVoteAction(action, username, postId)
	if err != nil {
		if err.Error() == ErrUnableToVoteUserAlreadyVoted {
			c.JSON(http.StatusOK, gin.H{"error": "user already voted"})
			return
		}
		s.logger.Errorf("unable to vote post: %v", err)
		c.JSON(500, gin.H{"error": "unable to vote post"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "voted"})
}

func (s *Server) routeUpvotePost(c *gin.Context) {
	s.innerRouteVotePost(c, models.ActionUpvote)
}

func (s *Server) routeDownvotePost(c *gin.Context) {
	s.innerRouteVotePost(c, models.ActionDownvote)
}
