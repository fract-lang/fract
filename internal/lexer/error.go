/*
	Error Function.
*/

package lexer

import (
	"fmt"
	"strings"

	"github.com/fract-lang/fract/pkg/str"
)

// Error Exit with error.
// message Message of error.
func (l Lexer) Error(message string) {
	fmt.Printf("File: %s\nPosition: %d:%d\n", l.File.Path, l.Line, l.Column)
	if !l.RangeComment { // Ignore multiline comment error.
		fmt.Println("    " + strings.ReplaceAll(strings.TrimSpace(l.File.Lines[l.Line-1]), "\t", " "))
		fmt.Println(str.GetWhitespace(4+l.Column-2) + "^")
	}
	fmt.Println(message)
	panic(nil)
}
