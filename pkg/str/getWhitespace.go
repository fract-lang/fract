package str

import "strings"

// GetWhitespace returns whitespace string by length.
func GetWhitespace(length int) string {
	var sb strings.Builder
	for length >= 0 {
		sb.WriteByte(' ')
		length--
	}
	return sb.String()
}
