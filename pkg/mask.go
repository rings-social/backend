package server

import "backend/pkg/models"

// maskDeletedComments
func maskDeletedComments(comments []models.Comment) []models.Comment {
	var maskedComments []models.Comment
	for _, v := range comments {
		if v.DeletedAt != nil {
			v.Body = "[deleted]"
		}
		maskedComments = append(maskedComments, v)
	}
	return maskedComments
}
