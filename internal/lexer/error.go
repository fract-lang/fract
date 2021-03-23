/*
	Error Function.
*/

package lexer

import (
	"fmt"
	"os"

	"github.com/fract-lang/fract/pkg/except"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/str"
)

// Error Exit with error.
// message Message of error.
func (l Lexer) Error(message string) {
	fmt.Printf("File: %s\nPosition: %d:%d\n",
		l.File.Path, l.Line, l.Column)
	if !l.MultilineComment { // Ignore multiline comment error.
		fmt.Println("    " + l.File.Lines[l.Line-1].Text)
		fmt.Println(str.GetWhitespace(4+l.Column-1) + "^")
	}
	fmt.Println(message)
	if l.File.Path == fract.Stdin {
		except.Raise("")
	}
	os.Exit(1)
}
