/*
	ReadyLines Function.
*/

package interpreter

import (
	"strings"
)

// ReadyLines Ready lines to process.
// lines Lines to ready.
func ReadyLines(lines []string) []string {
	readyLines := make([]string, len(lines))
	for index, line := range lines {
		readyLines[index] = strings.TrimRight(line, " \t\n\r")
	}
	return readyLines
}
