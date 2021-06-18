package lexer

import (
	"fmt"
	"strings"

	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/str"
)

// Lexer of Fract.
type Lexer struct {
	lastToken objects.Token

	File           *objects.SourceFile
	Column         int // Last column.
	Line           int // Last line.
	Finished       bool
	RangeComment   bool
	BraceCount     int
	BracketCount   int
	ParenthesCount int
}

// Error thrown exception.
func (l Lexer) Error(message string) {
	fmt.Printf("File: %s\nPosition: %d:%d\n", l.File.Path, l.Line, l.Column)
	if !l.RangeComment { // Ignore multiline comment error.
		fmt.Println("    " + strings.ReplaceAll(l.File.Lines[l.Line-1], "\t", " "))
		fmt.Println(str.Whitespace(4+l.Column-2) + "^")
	}
	fmt.Println(message)
	panic(nil)
}
