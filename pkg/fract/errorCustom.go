/*
	ErrorCustom Function.
*/

package fract

import (
	"fmt"
	"os"

	"github.com/fract-lang/fract/pkg/except"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/str"
)

// ErrorCustom Exit with error.
// file Code file instance.
// line Line of error.
// column Column of error.
// message Message of error.
func ErrorCustom(file obj.CodeFile, line, column int, message string) {
	fmt.Printf("File: %s\nPosition: %d:%d\n", file.Path, line, column)
	fmt.Println("    " + file.Lines[line-1].Text)
	fmt.Println(str.GetWhitespace(4+column-1) + "^")
	fmt.Println(message)
	if file.Path == Stdin {
		except.Raise("")
	}
	os.Exit(1)
}
