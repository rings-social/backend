package server

import (
	"backend/pkg/models"
	"fmt"
	"time"
)

func (s *Server) repoComments(postId uint, parentId *uint) ([]models.Comment, error) {
	var comments []models.Comment

	tx := s.db.
		Limit(200).
		Preload("Author").Order("score desc")
	var err error
	if parentId == nil {
		// Postgres doesn't like to compare NULLs with =, so we have to do this.
		err = tx.
			Find(&comments, "post_id = ? AND parent_id IS NULL", postId).Error
	} else {
		err = tx.
			Find(&comments, "post_id = ? AND parent_id = ?", postId, parentId).Error
	}

	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *Server) repoComment(commentId uint) (models.Comment, error) {
	var comment models.Comment
	tx := s.db.Preload("Author").First(&comment, "id = ?", commentId)
	return comment, tx.Error
}

func (s *Server) repoGetCommentTree(postId uint, commentId uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := s.db.Raw(`
		WITH RECURSIVE comment_tree AS (
			SELECT
				c.*
			FROM
				comments c
			WHERE
				c.id = ? AND post_id = ?
			UNION ALL
			SELECT
				cr.*
			FROM
				comments cr
				JOIN comment_tree ct ON cr.parent_id = ct.id
		)
		SELECT
			ct.*,
			u.username AS author_username
		FROM
			comment_tree ct
			JOIN users u ON ct.author_username = u.username
		ORDER BY ct.depth`,
		commentId,
		postId,
	).
		Scan(&comments).
		Error

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *Server) repoGetTopComments(postID uint) ([]models.Comment, error) {
	var comments []models.Comment

	maxLevel := 5
	topComments := 10
	err := s.db.Model(&models.Comment{}).Preload("Author").Raw(`
		WITH RECURSIVE comment_hierarchy AS (
			SELECT
				c.*,
				0 AS level
			FROM
				comments c
			WHERE
				c.post_id = ?
				AND c.parent_id IS NULL
			UNION ALL
			SELECT
				cr.*,
				ch.level + 1 AS level
			FROM
				comments cr
				JOIN comment_hierarchy ch ON cr.parent_id = ch.id
			WHERE
				ch.level < ?
		)
		SELECT
			ch.*,
			u.username AS author_username
		FROM
			comment_hierarchy ch
			JOIN users u ON ch.author_username = u.username
		ORDER BY
			ch.score DESC
		LIMIT ?
	`, postID, maxLevel, topComments).Scan(&comments).Error
	if err != nil {
		return nil, err
	}

	// Fill the users
	usernamesMap := map[string]bool{}
	for _, comment := range comments {
		usernamesMap[comment.AuthorUsername] = true
	}

	var usernames []string
	for username := range usernamesMap {
		usernames = append(usernames, username)
	}

	usersMap := map[string]models.User{}
	users, err := s.repoUsers(usernames)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		usersMap[user.Username] = user
	}

	for k, comment := range comments {
		// Check if user exists in map
		v, ok := usersMap[comment.AuthorUsername]
		if !ok {
			return nil, fmt.Errorf("user %s not found", comment.AuthorUsername)
		}
		comments[k].Author = v
	}

	return comments, nil
}

func (s *Server) repoDeleteComment(commentId uint) error {
	tx := s.db.Model(&models.Comment{}).
		Where("id = ?", commentId).
		Update("deleted_at", time.Now())
	return tx.Error
}
