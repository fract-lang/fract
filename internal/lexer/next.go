package lexer

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// Next lex next line.
func (l *Lexer) Next() []objects.Token {
	var tokens []objects.Token

	// If file is finished?
	if l.Finished {
		return tokens
	}

	// Reset bracket counter.
	l.ParenthesCount = 0
	l.BraceCount = 0
	l.BracketCount = 0

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
	for token.Type != fract.TypeNone &&
		token.Type != fract.TypeStatementTerminator {
		if !l.RangeComment && token.Type != fract.TypeIgnore {
			tokens = append(tokens, token)
			l.lastToken = token
		}
		token = l.Generate()
	}

	l.lastToken = token

	// Go next line.
	l.Line++

	// Line equals to or bigger then last line.
	l.Finished = l.Line > len(l.File.Lines)

	if l.ParenthesCount > 0 { // Check parentheses.
		if l.Finished {
			if l.File.Path != fract.Stdin {
				l.Line-- // Subtract for correct line number.
				l.Error("Parentheses is expected to close...")
			}
		} else {
			goto tokenize
		}
	} else if l.BraceCount > 0 { // Check braces.
		if l.Finished {
			if l.File.Path != fract.Stdin {
				l.Line-- // Subtract for correct line number.
				l.Error("Brace is expected to close...")
			}
		} else {
			goto tokenize
		}
	} else if l.BracketCount > 0 { // Check brackets.
		if l.Finished {
			if l.File.Path != fract.Stdin {
				l.Line-- // Subrract for correct line number.
				l.Error("Bracket is expected to close...")
			}
		} else {
			goto tokenize
		}
	} else if l.RangeComment {
		if l.Finished {
			if l.File.Path != fract.Stdin {
				l.Line--
				l.Error("Multiline comment is expected to close...")
			}
		} else {
			goto tokenize
		}
	}

	return tokens
}
