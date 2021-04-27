package str

import "strings"

// Returns string whitespace by count.
// count Count of whitespace.
func GetWhitespace(count int) string {
	var sb strings.Builder
	for count >= 0 {
		sb.WriteString(" ")
		count--
	}
	return sb.String()
}
