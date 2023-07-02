package server

import "backend/pkg/models"

// cleanupUsers removes sensitive information from a list of users.
func cleanupUsers(users []models.User) []models.User {
	for i := range users {
		users[i].AuthSubject = nil
	}
	return users
}
