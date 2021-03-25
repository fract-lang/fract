/*
	Error Function.
*/

package lexer

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/str"
)

// Error Exit with error.
// message Message of error.
func (l Lexer) Error(message string) {
	fmt.Printf("File: %s\nPosition: %d:%d\n",
		l.File.Path, l.Line, l.Column)
	if !l.RangeComment { // Ignore multiline comment error.
		fmt.Println("    " + l.File.Lines[l.Line-1].Text)
		fmt.Println(str.GetWhitespace(4+l.Column-1) + "^")
	}
	fmt.Println(message)
	panic(nil)
}
