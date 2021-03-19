package str

import "strings"

// Returns string whitespace by count.
// count Count of whitespace.
func GetWhitespace(count int) string {
	var sb strings.Builder
	for counter := 1; counter <= count; counter++ {
		sb.WriteRune(' ')
	}
	return sb.String()
}
