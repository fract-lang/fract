package str

import "strings"

// Whitespace returns whitespace string by length.
func Whitespace(len int) string {
	var sb strings.Builder
	for len >= 0 {
		sb.WriteByte(' ')
		len--
	}
	return sb.String()
}
