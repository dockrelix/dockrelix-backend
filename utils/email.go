package utils

import "regexp"

func IsEmailValid(email string) bool {
	regex := `^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`
	return regexp.MustCompile(regex).MatchString(email)
}
