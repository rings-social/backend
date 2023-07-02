package server

import (
	"backend/pkg/models"
	"fmt"
)

func (s *Server) repoGetTotalRings() (int64, error) {
	var total int64
	tx := s.db.Model(&models.Ring{}).Where("deleted_at IS NULL").Count(&total)
	if tx.Error != nil {
		s.logger.Errorf("Unable to get total rings: %v", tx.Error)
		return 0, fmt.Errorf("unable to get total rings")
	}

	return total, nil
}

func (s *Server) repoGetRings(offset int, limit uint64) ([]models.Ring, error) {
	var rings []models.Ring

	tx := s.db
	if offset != 0 {
		tx = tx.Offset(offset)
	}

	tx.Preload("Owner").
		Limit(int(limit)).
		Where("deleted_at IS NULL").
		Order("subscribers DESC").
		Find(&rings)
	if tx.Error != nil {
		s.logger.Errorf("Unable to get rings: %v", tx.Error)
		return nil, fmt.Errorf("unable to get rings")
	}

	return rings, nil
}
