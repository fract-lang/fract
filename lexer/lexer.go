package lexer

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"../fract"
	"../grammar"
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

// Last putted token.
var lastToken objects.Token

// Generate Generate next token.
func (l *Lexer) Generate() objects.Token {
	var (
		token objects.Token
		ln    string = l.File.Lines[l.Line-1].Text
	)

	/* Line is finished. */
	if l.Column >= len(ln) {
		return token
	}

	// Resume.
	ln = ln[l.Column-1:]

	/* Skip spaces. */
	for index := 0; index < len(ln); index++ {
		var char byte = ln[index]
		if char == ' ' || char == '\t' {
			l.Column++
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

	var arithmeticCheck string = strings.TrimSpace(regexp.MustCompile(
		"^(-|)\\s*[0-9]+(\\.[0-9]+)?(\\s+|$)").FindString(ln))
	if arithmeticCheck != "" { // Numeric value.
		token.Value = arithmeticCheck
		token.Type = fract.TypeValue
	} else if strings.HasPrefix(ln, grammar.TokenPlus) { // Addition.
		token.Value = grammar.TokenPlus
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.TokenMinus) { // Subtraction.
		token.Value = grammar.TokenMinus
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.TokenStar) { // Multiplication.
		token.Value = grammar.TokenStar
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.TokenSlash) { // Division.
		token.Value = grammar.TokenSlash
		token.Type = fract.TypeOperator
	} else {
		l.Error("What is this?: " + ln)
	}

	var tokenvallen int = len(token.Value)
	if tokenvallen == 1 {
		l.Column++
	} else {
		l.Column += tokenvallen - 1
	}

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
	lastToken.Type = fract.TypeNone
	lastToken.Line = 0
	lastToken.Column = 0
	lastToken.Value = ""

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
