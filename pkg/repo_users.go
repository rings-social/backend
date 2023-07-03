package server

import "backend/pkg/models"

// repoUsers returns a list of users matching the given usernames.
func (s *Server) repoUsers(usernames []string) ([]models.User, error) {
	var users []models.User
	err := s.db.
		Where("username IN (?)", usernames).
		Where("deleted_at IS NULL").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}
