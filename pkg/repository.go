package server

import (
	"backend/pkg/models"
	"gorm.io/gorm/clause"
)

func (s *Server) repoRingPosts(ringName string) ([]models.Post, error) {
	var ring models.Ring
	err := s.db.Preload(clause.Associations).First(&ring, "name = ?", ringName).Error
	if err != nil {
		return nil, err
	}

	return ring.Posts, nil
}
