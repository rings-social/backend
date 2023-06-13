package server

import "backend/pkg/models"

func (s *Server) repoRingsSearch(q string, nsfw bool) ([]models.Ring, error) {
	var rings []models.Ring
	tx := s.db.Find(&rings, "name LIKE ? AND nsfw IN (?, false)", "%"+q+"%", nsfw)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return rings, nil
}
