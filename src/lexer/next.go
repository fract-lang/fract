/*
	Next Function
*/

package lexer

import (
	"github.com/fract-lang/fract/src/fract"
	"github.com/fract-lang/fract/src/utils/vector"
)

// Next Lex next line.
func (l *Lexer) Next() vector.Vector {
	tokens := vector.Vector{}

	// If file is finished?
	if l.Finished {
		return tokens
	}

	// Reset bracket counter.
	l.parenthesCount = 0
	l.braceCount = 0
	l.bracketCount = 0

tokenize:

	if l.lastToken.Type != fract.TypeStatementTerminator {
		// Restore to defaults.
		l.Column = 1
		l.lastToken.Type = fract.TypeNone
		l.lastToken.Line = 0
		l.lastToken.Column = 0
		l.lastToken.Value = ""
	}

	// Tokenize line.
	token := l.Generate()
	for token.Type != fract.TypeNone && token.Type != fract.TypeStatementTerminator {
		tokens.Vals = append(tokens.Vals, token)
		l.lastToken = token
		token = l.Generate()
	}

	l.lastToken = token

	// Go next line.
	l.Line++

	// Line equals to or bigger then last line.
	l.Finished = l.Line > len(l.File.Lines.Vals)

	if l.parenthesCount > 0 { // Check parentheses.
		if l.Finished {
			l.Line-- // Subtract for correct line number.
			l.Error("Parentheses is expected to close...")
		}
		goto tokenize
	} else if l.braceCount > 0 { // Check braces.
		if l.Finished {
			l.Line-- // Subtract for correct line number.
			l.Error("Brace is expected to close...")
		}
		goto tokenize
	} else if l.bracketCount > 0 { // Check brackets.
		if l.Finished {
			l.Line-- // Subrract for correct line number.
			l.Error("Bracket is expected to close...")
		}
		goto tokenize
	}

	return tokens
}
