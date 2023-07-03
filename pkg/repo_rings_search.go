package server

import (
	"backend/pkg/models"
	"strconv"
	"strings"
)

func (s *Server) repoRingsSearch(q string, nsfw bool, after string, limit uint64) ([]models.Ring, error) {
	var rings []models.Ring
	q = strings.Replace(q, "%", "\\%", -1)
	q = strings.ToLower(q)

	tx := s.db
	if after != "" {
		afterV, err := strconv.ParseUint(after, 10, 64)
		if err != nil {
			return nil, err
		}
		tx = tx.Where("subscribers < ?", afterV)
	}
	tx.
		Preload("Owner").
		Limit(int(limit)).
		Order("subscribers DESC").
		Find(&rings, "name LIKE ? AND nsfw IN (?, false)", q+"%", nsfw)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return rings, nil
}
