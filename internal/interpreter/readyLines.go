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
	var readyLines []string
	for index := 0; index < len(lines); index++ {
		readyLines = append(readyLines, strings.TrimRight(lines[index], " \t\n\r"))
	}
	return readyLines
}
