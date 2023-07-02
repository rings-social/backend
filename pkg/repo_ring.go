package server

import "backend/pkg/models"

func (s *Server) repoRingExists(ringName string) (bool, error) {
	var ring models.Ring
	err := s.db.First(&ring, "name = ?", ringName).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
