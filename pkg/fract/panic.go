package fract

import (
	"fmt"
	"strings"

	"github.com/fract-lang/fract/pkg/obj"
	"github.com/fract-lang/fract/pkg/str"
)

func PanicC(f *obj.File, col, ln int, t, m string) {
	e := obj.Panic{
		M: fmt.Sprintf("File: %s\nPosition: %d:%d\n    %s\n%s^\n%s: %s",
			f.P, ln, col, strings.ReplaceAll(f.Lns[ln-1], "\t", " "),
			str.Whitespace(4+col-2), t, m),
		T: t,
	}
	if TryCount > 0 {
		panic(e)
	}
	e.Panic()
}

func Panic(tk obj.Token, t, m string) { PanicC(tk.F, tk.Col, tk.Ln, t, m) }

// Interpreter panic.
func IPanicC(f *obj.File, ln, col int, t, m string) {
	e := obj.Panic{
		M: fmt.Sprintf("File: %s\nPosition: %d:%d\n    %s\n%s^\n%s: %s",
			f.P, ln, col, strings.ReplaceAll(f.Lns[ln-1], "\t", " "),
			str.Whitespace(4+col-2), t, m),
		T: t,
	}
	e.Panic()
}

// Interpreter panic.
func IPanic(tk obj.Token, t, m string) { IPanicC(tk.F, tk.Ln, tk.Col, t, m) }
