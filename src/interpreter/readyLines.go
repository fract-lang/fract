/*
	ReadyLines Function.
*/

package interpreter

import (
	"strings"

	"github.com/fract-lang/fract/src/objects"
)

// ReadyLines Ready lines to process.
// lines Lines to ready.
func ReadyLines(lines []string) []objects.CodeLine {
	readyLines := make([]objects.CodeLine, 0)
	for index := 0; index < len(lines); index++ {
		readyLines = append(readyLines, objects.CodeLine{Line: index + 1,
			Text: strings.TrimRight(lines[index], " \t\n\r")})
	}
	return readyLines
}
