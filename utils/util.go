package utils

import "regexp"

var emailRe = regexp.MustCompile(`^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*(\.[a-zA-Z]{2,6})$`)

func IsValidEmail(email string) bool {
	return emailRe.MatchString(email)
}
