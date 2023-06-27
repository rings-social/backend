package server

import (
	"backend/pkg/request"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
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

	comments, done := s.retrieveComments(c, postId, parentId)
	if done {
		return
	}

	comments = maskDeletedComments(comments)

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
