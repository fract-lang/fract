/*
	ReadyLines Function.
*/

package interpreter

import (
	"strings"

	obj "github.com/fract-lang/fract/pkg/objects"
)

// ReadyLines Ready lines to process.
// lines Lines to ready.
func ReadyLines(lines []string) []obj.CodeLine {
	readyLines := make([]obj.CodeLine, 0)
	for index := 0; index < len(lines); index++ {
		readyLines = append(readyLines, obj.CodeLine{Line: index + 1,
			Text: strings.TrimRight(lines[index], " \t\n\r")})
	}
	return readyLines
}
