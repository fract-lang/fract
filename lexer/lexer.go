package lexer

import (
	"../fract"
	"../objects"
	"../utilities/list"
)

// Lexer of Fract.
type Lexer struct {
	/* PUBLIC */

	// Destination file.
	File *objects.CodeFile
	// Last column.
	Column int
	// Last line.
	Line int
	// Finished file.
	Finished bool
}

// New Create new instance.
func New(file objects.CodeFile) *Lexer {
	var lexer *Lexer = new(Lexer)
	lexer.File = &file
	lexer.Line = 1
	return lexer
}

// Generate Generate next token.
func (l *Lexer) Generate() objects.Token {
	var (
		token objects.Token
		ln    string = l.File.Lines[l.Line-1].Text
	)

	if ln == "" || l.Column >= len(ln) {
		return token
	}

	token.Column = l.Column
	token.Line = l.Line
	token.Type = fract.TypePrint
	token.Value = ln
	l.Column += len(token.Value)

	return token
}

// Next Lex next line.
func (l *Lexer) Next() list.List {
	var tokens *list.List = list.New()

	// If file is finished?
	if l.Finished {
		return *tokens
	}

	// Restore to defaults.
	l.Column = 1

	// Tokenize line.
	var token objects.Token = l.Generate()
	for token.Value != "" {
		tokens.Append(token)
		token = l.Generate()
	}

	// Go next line.
	l.Line++
	// Line equals to or biggren then last line.
	l.Finished = l.Line > len(l.File.Lines)

	return *tokens
}
