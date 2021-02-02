/*
	Next Function
*/

package lexer

import (
	"../fract"
	"../utilities/vector"
)

// Next Lex next line.
func (l *Lexer) Next() vector.Vector {
	tokens := vector.New()

	// If file is finished?
	if l.Finished {
		return *tokens
	}

	// Restore to defaults.
	l.Column = 1
	l.lastToken.Type = fract.TypeNone
	l.lastToken.Line = 0
	l.lastToken.Column = 0
	l.lastToken.Value = ""
	l.braceCount = 0

	// Tokenize line.
	token := l.Generate()
	for token.Value != "" {
		tokens.Append(token)
		l.lastToken = token
		token = l.Generate()
	}

	/* Check parentheses. */
	if l.braceCount > 0 {
		l.Error("Bracket is expected to close...")
	}

	// Go next line.
	l.Line++

	// Line equals to or bigger then last line.
	l.Finished = l.Line > len(l.File.Lines.Vals)

	return *tokens
}
