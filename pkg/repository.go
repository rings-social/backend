package server

import (
	"backend/pkg/models"
	"gorm.io/gorm/clause"
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

func (s *Server) repoGetUserByAuthSubject(subject string) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, "auth_subject = ?", subject).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
