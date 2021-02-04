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

	/* Check arithmetic value? */
	check := strings.TrimSpace(regexp.MustCompile(
		"^(-|)\\s*[0-9]+(\\.[0-9]+)?(\\s+||\\W|$)").FindString(ln))
	if check != "" &&
		(l.lastToken.Value == "" || l.lastToken.Type == fract.TypeOperator ||
			l.lastToken.Type == fract.TypeBrace) { // Numeric value.
		match, _ := regexp.MatchString("\\W$", check)
		if match {
			check = check[:len(check)-1]
		}
		token.Value = check
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
	} else if strings.HasPrefix(ln, grammar.Setter) {
		token.Value = grammar.Setter
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.KwVariable) { // Variable.
		token.Value = grammar.KwVariable
		token.Type = fract.TypeVariable
	} else if strings.HasPrefix(ln, grammar.DtInt8) { // int8.
		token.Value = grammar.DtInt8
		token.Type = fract.TypeDataType
	} else if strings.HasPrefix(ln, grammar.DtInt16) { // int16.
		token.Value = grammar.DtInt16
		token.Type = fract.TypeDataType
	} else if strings.HasPrefix(ln, grammar.DtInt32) { // int32.
		token.Value = grammar.DtInt32
		token.Type = fract.TypeDataType
	} else if strings.HasPrefix(ln, grammar.DtInt64) { // int64.
		token.Value = grammar.DtInt64
		token.Type = fract.TypeDataType
	} else if strings.HasPrefix(ln, grammar.DtUInt8) { // uint8.
		token.Value = grammar.DtUInt8
		token.Type = fract.TypeDataType
	} else if strings.HasPrefix(ln, grammar.DtUInt16) { // uint16.
		token.Value = grammar.DtUInt16
		token.Type = fract.TypeDataType
	} else if strings.HasPrefix(ln, grammar.DtUInt32) { // uint32.
		token.Value = grammar.DtUInt32
		token.Type = fract.TypeDataType
	} else if strings.HasPrefix(ln, grammar.DtUInt64) { // uint64.
		token.Value = grammar.DtUInt64
		token.Type = fract.TypeDataType
	} else if strings.HasPrefix(ln, grammar.DtFloat32) { // float32.
		token.Value = grammar.DtFloat32
		token.Type = fract.TypeDataType
	} else if strings.HasPrefix(ln, grammar.DtFloat64) { // float64.
		token.Value = grammar.DtFloat64
		token.Type = fract.TypeDataType
	} else if strings.HasPrefix(ln, grammar.TokenSharp) { // Comment.
	} else { // Alternates
		/* Check variable name. */
		check = strings.TrimSpace(regexp.MustCompile(
			"^([A-z])([a-zA-Z1-9" + grammar.TokenUnderscore + grammar.TokenDot +
				".]+)?(\\s+|$)").FindString(ln))
		if check != "" && !strings.HasSuffix(check, grammar.TokenDot) &&
			!strings.HasSuffix(check, grammar.TokenUnderscore) { // Name.
			token.Value = strings.TrimSpace(check)
			token.Type = fract.TypeName
		} else { // Error exactly
			l.Error("What is this?: " + ln)
		}

	}

	/* Add length to column. */
	l.Column += len(token.Value)

	return token
}
