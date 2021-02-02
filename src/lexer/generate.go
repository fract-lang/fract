/*
	Generate Function
*/

package lexer

import (
	"regexp"
	"strings"

	"../fract"
	"../grammar"
	"../objects"
)

// Generate Generate next token.
func (l *Lexer) Generate() objects.Token {
	var token objects.Token
	ln := l.File.Lines.At(l.Line - 1).(objects.CodeLine).Text

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
		(l.lastToken.Value == "" || l.lastToken.Type == fract.TypeOperator ||
			l.lastToken.Type == fract.TypeBrace) { // Numeric value.
		match, _ := regexp.MatchString("\\W$", arithmeticCheck)
		if match {
			arithmeticCheck = arithmeticCheck[:len(arithmeticCheck)-1]
		}
		token.Value = arithmeticCheck
		token.Type = fract.TypeValue
	} else if strings.HasPrefix(ln, grammar.IntegerDivision) { // Integer division.
		token.Value = grammar.IntegerDivision
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.IntegerDivideWithBigger) { // Integer divide with bigger.
		token.Value = grammar.IntegerDivideWithBigger
		token.Type = fract.TypeOperator
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
	} else if strings.HasPrefix(ln, grammar.TokenCaret) { // Exponentiation.
		token.Value = grammar.TokenCaret
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.TokenPercent) { // Mod.
		token.Value = grammar.TokenPercent
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.TokenReverseSlash) { // Divisin with bigger.
		token.Value = grammar.TokenReverseSlash
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.TokenLParenthes) { // Open parentheses.
		l.braceCount++
		token.Value = grammar.TokenLParenthes
		token.Type = fract.TypeBrace
	} else if strings.HasPrefix(ln, grammar.TokenRParenthes) { // Close parentheses.
		l.braceCount--
		if l.braceCount < 0 {
			l.Error("The extra parentheses are closed!")
		}
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
