package server

import "backend/pkg/models"

func (s *Server) repoGetVote(username string, postId uint) (*models.PostAction, error) {
	var postAction models.PostAction
	tx := s.db.First(&postAction, "username = ? AND post_id = ?", username, postId)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &postAction, nil
}
