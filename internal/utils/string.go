package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func FindAllInteger(msg string) int {
	// Compile a regular expression to match one or more digits (\\d+)
	re := regexp.MustCompile(`\d+`)
	// Find all strings that match the pattern, -1 means find all matches
	matches := strings.Join(re.FindAllString(msg, -1), "")
	res, _ := strconv.Atoi(matches)

	return res
}
