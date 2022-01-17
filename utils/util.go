package utils

import "regexp"

var emailRe = regexp.MustCompile(`[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z]{2,6}`)

func NormalEmail(email string) string {
	if v := emailRe.FindStringSubmatch(email); len(v) > 0 {
		return v[0]
	}

	return ""
}