package fract

import (
	"fmt"
	"strings"

	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/str"
)

// ErrorCustom thrown new exception.
func ErrorCustom(file *objects.SourceFile, line, column int, message string) {
	e := objects.Exception{
		Message: fmt.Sprintf("File: %s\nPosition: %d:%d\n    %s\n%s^\n%s",
			file.Path, line, column, strings.ReplaceAll(file.Lines[line-1], "\t", " "),
			str.GetWhitespace(4+column-2), message),
	}
	if TryCount > 0 {
		panic(fmt.Errorf(e.Message))
	}
	fmt.Println(e.Message)
	panic(fmt.Errorf(""))
}
