/*
	ErrorCustom Function.
*/

package fract

import (
	"fmt"
	"strings"

	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/str"
)

// ErrorCustom Create new exception.
// file Code file instance.
// line Line of error.
// column Column of error.
// message Message of error.
func ErrorCustom(file obj.CodeFile, line, column int, message string) obj.Exception {
	e := obj.Exception{
		Message: fmt.Sprintf("File: %s\nPosition: %d:%d\n    %s\n%s^\n%s",
			file.Path, line, column, strings.ReplaceAll(strings.TrimSpace(file.Lines[line-1]), "\t", " "),
			str.GetWhitespace(4+column-1), message),
	}

	if TryCount > 0 {
		panic(fmt.Errorf(e.Message))
	}

	fmt.Println(e.Message)
	panic(nil)
}
