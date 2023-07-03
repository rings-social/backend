package server

import (
	"backend/pkg/models"
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

	comments, done := s.retrieveComments(c, uint(postId), parentId, map[uint]models.CommentAction{})
	if done {
		return
	}

	comments = maskDeletedComments(comments)

	redditComments, err := toRedditComments(post, comments, s.baseUrl)
	if err != nil {
		internalServerError(c)
		return
	}
	c.JSON(200, redditComments)
}

func (s *Server) retrieveComments(
	c *gin.Context,
	postId uint,
	parentId *uint,
	commentActions map[uint]models.CommentAction,
) ([]models.Comment, bool) {

	var comments []models.Comment
	var err error
	if parentId != nil {
		comments, err = s.repoGetCommentTree(postId, *parentId)
	} else {
		comments, err = s.repoGetTopComments(postId)
	}
	if err != nil {
		s.logger.Errorf("unable to get comments for post %d: %v", postId, err)
		internalServerError(c)
		return nil, true
	}
	comments = setCommentActions(comments, commentActions)
	comments = maskDeletedComments(comments)

	// Create a tree structure
	commentTree := map[uint][]models.Comment{}
	var topLevelComments []models.Comment
	for k, comment := range comments {
		if comment.ParentId != nil {
			commentTree[*comment.ParentId] = append(commentTree[*comment.ParentId], comments[k])
		} else {
			topLevelComments = append(topLevelComments, comments[k])
		}
	}

	// Now that we have the relationships, create the array of top level comments (no parent)
	for k, comment := range topLevelComments {
		topLevelComments[k] = fillChildren(1, &comment, commentTree)
	}

	return topLevelComments, false
}

// fillChildren recursively fills the children with a comment
func fillChildren(depth int, c *models.Comment, tree map[uint][]models.Comment) models.Comment {
	c.Depth = depth
	children, ok := tree[c.ID]
	if !ok {
		return *c
	}
	for k, child := range children {
		children[k] = fillChildren(depth+1, &child, tree)
	}
	c.Replies = children
	return *c
}

func setCommentActions(comments []models.Comment, actions map[uint]models.CommentAction) []models.Comment {
	for k, v := range comments {
		action, ok := actions[v.ID]
		if ok {
			switch action.Action {
			case models.ActionDownvote:
				comments[k].VotedDown = true
			case models.ActionUpvote:
				comments[k].VotedUp = true
			}
		}
	}
	return comments
}

func setDepth(comments []models.Comment, i int) {
	for k, _ := range comments {
		comments[k].Depth = i
	}
}
