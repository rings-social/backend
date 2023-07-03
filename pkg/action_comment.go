package server

import (
	"backend/pkg/models"
	"backend/pkg/request"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (s *Server) getComments(c *gin.Context) {
	postId, done := parsePostId(c)
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

	commentActions := map[uint]models.CommentAction{}
	if s.hasIdToken(c) {
		idToken, done := s.idToken(c)
		if done {
			return
		}

		username, err := s.usernameForIdToken(idToken)
		if err == nil {
			commentActionsArr := s.repoCommentActions(username, postId)
			for _, v := range commentActionsArr {
				commentActions[v.CommentId] = v
			}
		}
	}

	comments, done := s.retrieveComments(c, uint(postId), parentId, commentActions)
	if done {
		return
	}

	c.JSON(200, comments)
}
func (s *Server) postComment(c *gin.Context) {
	postId, done := parsePostId(c)
	if done {
		return
	}

	// Check if user is authenticated
	idToken, done := s.idToken(c)
	if done {
		return
	}

	// Get user id by idToken
	username, err := s.usernameForIdToken(idToken)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "invalid id token",
		})
		return
	}

	var commentRequest request.Comment
	err = c.BindJSON(&commentRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	comment, err := s.addComment(uint(postId), username, commentRequest)
	if err != nil {
		s.logger.Errorf("unable to add comment: %v", err)
		internalServerError(c)
		return
	}

	c.JSON(200, comment)
}

func (s *Server) deleteComment(c *gin.Context) {
	postId, done := parsePostId(c)
	if done {
		return
	}

	commentId, done := parseId(c, "commentId")
	if done {
		return
	}

	// Check if user is authenticated
	idToken, done := s.idToken(c)
	if done {
		return
	}

	// Get user by idtoken
	username, err := s.usernameForIdToken(idToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user not registered",
		})
		return
	}

	// Get comment
	comment, err := s.repoComment(uint(commentId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "comment not found",
			})
			return
		}
	}

	if comment.PostId != uint(postId) {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "comment doesn't belong to this post",
			},
		)
		return
	}

	if comment.AuthorUsername != username {
		// Check if the user is an admin
		if !s.isAdmin(username) {
			s.logger.Errorf("user %s tried to delete comment %d which he doesn't own", username, commentId)
			c.JSON(http.StatusForbidden,
				gin.H{"error": "you can't delete this comment"},
			)
			return
		}
	}

	if comment.DeletedAt != nil {
		// Comment already deleted
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "comment already deleted"})
		return
	}

	err = s.repoDeleteComment(uint(commentId))
	if err != nil {
		s.logger.Errorf("unable to delete comment: %v", err)
		internalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (s *Server) voteAction(c *gin.Context, action models.VoteAction) {
	// Needs to be logged in to proceed
	idToken, done := s.idToken(c)
	if done {
		return
	}
	username, err := s.usernameForIdToken(idToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user not registered",
		})
		return
	}

	// Get comment
	commentId, done := parseId(c, "commentId")
	if done {
		return
	}
	comment, err := s.repoComment(uint(commentId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "comment not found",
			})
			return
		}
		internalServerError(c)
		return
	}

	if comment.DeletedAt != nil {
		// Comment already deleted
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("deleted comments cannot be %s", action)})
		return
	}

	voteAction, addScore, err := s.repoCommentVoteAction(uint(commentId), username, action)
	if err != nil {
		// Unable to upvote, bad request:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("unable to %s comment", action),
		})
		return
	}

	err = s.repoIncreaseCommentScore(uint(commentId), addScore)
	if err != nil {
		s.logger.Errorf("unable to change comment score: %v", err)
		internalServerError(c)
		return
	}

	c.JSON(http.StatusOK, voteAction)
}

func (s *Server) upvoteComment(c *gin.Context) {
	s.voteAction(c, models.ActionUpvote)
}
func (s *Server) downvoteComment(c *gin.Context) {
	s.voteAction(c, models.ActionDownvote)
}
