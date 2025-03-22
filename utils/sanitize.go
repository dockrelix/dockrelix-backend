package utils

import (
	"strings"
)

func SanitizeInput(input string) string {
	replacer := strings.NewReplacer(
		"'", "", "\"", "", ";", "", "--", "", "/*", "", "*/", "", "xp_", "", "/", "",
	)
	return replacer.Replace(input)
}
