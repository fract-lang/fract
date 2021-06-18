package str

import "strings"

// Whitespace returns whitespace string by length.
func Whitespace(length int) string {
	var sb strings.Builder
	for length >= 0 {
		sb.WriteByte(' ')
		length--
	}
	return sb.String()
}
