package lexer

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"../fract"
	"../grammar"
	"../objects"
	"../utilities/vector"
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
	lexer := new(Lexer)
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
	var token objects.Token
	ln := l.File.Lines[l.Line-1].Text

	/* Line is finished. */
	if l.Column > len(ln) {
		return token
	}

	// Resume.
	ln = ln[l.Column-1:]

	/* Skip spaces. */
	for index := 0; index < len(ln); index++ {
		char := ln[index]
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

	/* Set token values. */
	token.Column = l.Column
	token.Line = l.Line

	/* Tokenize. */
	arithmeticCheck := strings.TrimSpace(regexp.MustCompile(
		"^(-|)\\s*[0-9]+(\\.[0-9]+)?(\\s+||\\W|$)").FindString(ln))
	if arithmeticCheck != "" &&
		(lastToken.Value == "" || lastToken.Type == fract.TypeOperator ||
			lastToken.Type == fract.TypeBrace) { // Numeric value.
		match, _ := regexp.MatchString("\\W$", arithmeticCheck)
		if match {
			arithmeticCheck = arithmeticCheck[:len(arithmeticCheck)-1]
		}
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
	} else if strings.HasPrefix(ln, grammar.TokenLParenthes) { // Open parentheses.
		token.Value = grammar.TokenLParenthes
		token.Type = fract.TypeBrace
	} else if strings.HasPrefix(ln, grammar.TokenRParenthes) { // Close parentheses.
		token.Value = grammar.TokenRParenthes
		token.Type = fract.TypeBrace
	} else if strings.HasPrefix(ln, grammar.TokenSharp) { // Comment.
	} else {
		l.Error("What is this?: " + ln)
	}

	/* Add length to column. */
	l.Column += len(token.Value)

	return token
}

// Next Lex next line.
func (l *Lexer) Next() vector.Vector {
	tokens := vector.New()

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
	token := l.Generate()
	for token.Value != "" {
		tokens.Append(token)
		lastToken = token
		token = l.Generate()
	}

	// Go next line.
	l.Line++

	// Line equals to or bigger then last line.
	l.Finished = l.Line > len(l.File.Lines)

	return *tokens
}
