package server

import (
	"backend/pkg/models"
	"gorm.io/gorm/clause"
	"time"
)

func (s *Server) repoRingPosts(ringName string) ([]models.Post, error) {
	var ring models.Ring
	err := s.db.
		Limit(100).
		Preload(clause.Associations).First(&ring, "name = ?", ringName).Error
	if err != nil {
		return nil, err
	}

	return ring.Posts, nil
}

func (s *Server) repoPost(postId uint) (*models.Post, error) {
	var post models.Post
	err := s.db.Preload(clause.Associations).First(&post, "id = ?", postId).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *Server) repoRingAbout(ringName string) (*models.Ring, error) {
	var ring models.Ring
	err := s.db.First(&ring, "name = ?", ringName).Error
	if err != nil {
		return nil, err
	}
	return &ring, nil
}

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

func (s *Server) repoDeleteComment(commentId uint) error {
	tx := s.db.Model(&models.Comment{}).
		Where("id = ?", commentId).
		Update("deleted_at", time.Now())
	return tx.Error
}

func (s *Server) repoGetUserByAuthSubject(subject string) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, "auth_subject = ?", subject).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
