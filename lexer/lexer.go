package lexer

import (
	"fmt"
	"os"
	"regexp"

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

// Error Exit with error.
// message Message of error.
func (l *Lexer) Error(message string) {
	fmt.Printf("LEXER ERROR\nMessage: %s\nLINE: %d\nCOLUMN: %d",
		message, l.Line, l.Column)
	os.Exit(1)
}

// Generate Generate next token.
func (l *Lexer) Generate() objects.Token {
	var (
		token objects.Token
		ln    string = l.File.Lines[l.Line-1].Text
	)

	/* Line is finished. */
	if l.Column > len(ln) {
		return token
	}

	// Resume.
	ln = ln[l.Column-1:]

	/* Skip spaces. */
	for index := 0; index < len(ln); index++ {
		l.Column++
		var char byte = ln[index]
		if char == ' ' || char == '\t' {
			continue
		}
		ln = ln[index:]
		break
	}

	/* Content is empty. */
	if ln == "" {
		return token
	}

	token.Column = l.Column
	token.Line = l.Line

	var arithmeticCheck string = regexp.MustCompile(
		"^(-|)\\s*[0-9]+(\\.[0-9]+)?(\\s+|$)").FindString(ln)
	if arithmeticCheck != "" {
		token.Value = arithmeticCheck
		token.Type = fract.TypeValue
	} else {
		l.Error("What is this?: " + ln)
	}

	l.Column += len(token.Value) - 1

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

	// Line equals to or bigger then last line.
	l.Finished = l.Line > len(l.File.Lines)

	return *tokens
}
