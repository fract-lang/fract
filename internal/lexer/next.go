package lexer

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// Check expected bracket or like and returns true if require retokenize, returns false if not.
// Thrown exception is syntax error.
func (l *Lexer) checkExpected(message string) bool {
	if l.Finished {
		if l.File.Path != "<stdin>" {
			l.Line-- // Subtract for correct line number.
			l.Error(message)
		}
		return false
	}
	return true
}

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

	switch {
	case l.ParenthesCount > 0: // Check parentheses.
		if l.checkExpected("Parentheses is expected to close...") { goto tokenize }
	case l.BraceCount > 0: // Check braces.
		if l.checkExpected("Brace is expected to close...") { goto tokenize }
	case l.BracketCount > 0: // Check brackets.
		if l.checkExpected("Bracket is expected to close...") { goto tokenize }
	case l.RangeComment:
		if l.checkExpected("Multiline comment is expected to close...") { goto tokenize}
	}

	return tokens
}
