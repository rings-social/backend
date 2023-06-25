package server

import "regexp"

func validateUsername(username string) (bool, string) {
	if len(username) < 4 {
		return false, "Username must be at least 4 characters long"
	}

	if len(username) > 20 {
		return false, "Username must be at most 20 characters long"
	}

	// Make sure it's only alphanumeric, numbers, and underscores
	m, err := regexp.MatchString("^[a-zA-Z0-9_]*$", username)
	if err != nil {
		return false, "Unable to validate username"
	}

	if !m {
		return false, "Username must only contain letters, numbers, and underscores"
	}

	invalidUsernames := []string{
		"admin",
		"system",
		"root",
		"moderator",
		"mod",
		"administrator",
		"me",
	}

	for _, invalidUsername := range invalidUsernames {
		if username == invalidUsername {
			return false, "Username is not allowed"
		}
	}
	return true, ""
}
