package fract

import (
	"fmt"
	"strings"

	"github.com/fract-lang/fract/pkg/obj"
	"github.com/fract-lang/fract/pkg/str"
)

// Errorc thrown new exception.
func Errorc(f *obj.File, ln, col int, msg string) {
	e := obj.Exception{
		Msg: fmt.Sprintf("File: %s\nPosition: %d:%d\n    %s\n%s^\n%s",
			f.P, ln, col, strings.ReplaceAll(f.Lns[ln-1], "\t", " "),
			str.Whitespace(4+col-2), msg),
	}
	if TryCount > 0 {
		panic(fmt.Errorf(e.Msg))
	}
	fmt.Println(e.Msg)
	panic(fmt.Errorf(""))
}

// Error thrown exception.
func Error(tk obj.Token, msg string) { Errorc(tk.F, tk.Ln, tk.Col, msg) }
