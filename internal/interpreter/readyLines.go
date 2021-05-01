package interpreter

import (
	"strings"
)

// ReadyLines returns lines processed to lexing.
func ReadyLines(lines []string) []string {
	readyLines := make([]string, len(lines))
	for index, line := range lines {
		readyLines[index] = strings.TrimRight(line, " \t\n\r")
	}
	return readyLines
}
