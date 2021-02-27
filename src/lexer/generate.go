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

// processEsacepeSequence Process char literal espace sequence.
// l Lexer.
// token Token.
// fln Full line text of current code line.
func processEscapeSequence(l *Lexer, token *objects.Token, fln string) bool {
	// Is not espace sequence?
	if fln[l.Column-1] != '\\' {
		return false
	}

	l.Column++
	if l.Column >= len(fln)+1 {
		l.Error("Charray literal is not defined full!")
	}

	switch fln[l.Column-1] {
	case '\\':
		token.Value += "\\"
	case '"':
		token.Value += "\""
	case '\'':
		token.Value += "'"
	case 'n':
		token.Value += "\n"
	case 'r':
		token.Value += "\r"
	case 't':
		token.Value += "\t"
	case 'b':
		token.Value += "\b"
	case 'f':
		token.Value += "\f"
	case 'a':
		token.Value += "\a"
	case 'v':
		token.Value += "\v"
	default:
		l.Error("Invalid escape sequence!")
	}

	return true
}

// lexChar Lex char literal.
// l Lexer.
// token Token.
// fln Full line text of current code line.
func lexChar(l *Lexer, token *objects.Token, fln string) {
	token.Value = grammar.TokenQuote
	token.Type = fract.TypeValue
	l.Column++
	for ; l.Column < len(fln)+1; l.Column++ {
		current := string(fln[l.Column-1])
		if current == grammar.TokenQuote { // Finish?
			token.Value += current
			break
		} else if !processEscapeSequence(l, token, fln) {
			token.Value += current
		}
	}
	if !strings.HasSuffix(token.Value, grammar.TokenQuote) {
		l.Error("Close quote is not found!")
	} else if len(token.Value) != 3 {
		l.Error("Char is only be one character!")
	}
	l.Column -= 2
}

// lexString Lex string literal.
// l Lexer.
// token Token.
// fln Full line text of current code line.
func lexString(l *Lexer, token *objects.Token, fln string) {
	token.Value = grammar.TokenDoubleQuote
	l.Column++
	for ; l.Column < len(fln)+1; l.Column++ {
		current := string(fln[l.Column-1])
		if current == grammar.TokenDoubleQuote { // Finish?
			token.Value += current
			break
		} else if !processEscapeSequence(l, token, fln) {
			token.Value += current
		}
	}
	if !strings.HasSuffix(token.Value, grammar.TokenDoubleQuote) {
		l.Error("Close double quote is not found!")
	}
	token.Type = fract.TypeValue

	l.Column -= len(token.Value) - 1
}

// Generate Generate next token.
func (l *Lexer) Generate() objects.Token {
	var token objects.Token
	token.File = l.File
	fln := l.File.Lines.Vals[l.Line-1].(objects.CodeLine).Text // Full line.

	/* Line is finished. */
	if l.Column > len(fln) {
		return token
	}

	// Resume.
	ln := fln[l.Column-1:]

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
	if check := strings.TrimSpace(regexp.MustCompile(
		"^(-|)\\s*[0-9]+(\\.[0-9]+)?(\\s|[[:punct:]]|$)").FindString(ln)); check != "" &&
		(l.lastToken.Value == "" || l.lastToken.Type == fract.TypeOperator ||
			(l.lastToken.Type == fract.TypeBrace && (l.lastToken.Value != grammar.TokenRBracket &&
				l.lastToken.Value != grammar.TokenRBracket)) || l.lastToken.Type == fract.TypeBlock ||
			l.lastToken.Type == fract.TypeStatementTerminator || l.lastToken.Type == fract.TypeLoop ||
			l.lastToken.Type == fract.TypeComma || l.lastToken.Type == fract.TypeIn ||
			l.lastToken.Type == fract.TypeIf || l.lastToken.Type == fract.TypeElseIf ||
			l.lastToken.Type == fract.TypeExit) { // Numeric value.
		// Remove punct.
		result, _ := regexp.MatchString("(\\s|[[:punct:]])$", check)
		if result {
			check = check[:len(check)-1]
		}
		clen := len(check)
		check = strings.ReplaceAll(check, " ", "")
		l.Column += clen - len(check)
		token.Value = check
		token.Type = fract.TypeValue
	} else if strings.HasPrefix(ln, grammar.TokenSemicolon) { // Statement terminator.
		token.Value = grammar.TokenSemicolon
		token.Type = fract.TypeStatementTerminator
		l.Line--
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
	} else if strings.HasPrefix(ln, grammar.TokenBackslash) { // Divisin with bigger.
		token.Value = grammar.TokenBackslash
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.TokenLParenthes) { // Open parentheses.
		l.parenthesCount++
		token.Value = grammar.TokenLParenthes
		token.Type = fract.TypeBrace
	} else if strings.HasPrefix(ln, grammar.TokenRParenthes) { // Close parentheses.
		l.parenthesCount--
		if l.parenthesCount < 0 {
			l.Error("The extra parentheses are closed!")
		}
		token.Value = grammar.TokenRParenthes
		token.Type = fract.TypeBrace
	} else if strings.HasPrefix(ln, grammar.TokenLBrace) { // Open brace.
		l.braceCount++
		token.Value = grammar.TokenLBrace
		token.Type = fract.TypeBrace
	} else if strings.HasPrefix(ln, grammar.TokenRBrace) { // Close brace.
		l.braceCount--
		if l.braceCount < 0 {
			l.Error("The extra brace are closed!")
		}
		token.Value = grammar.TokenRBrace
		token.Type = fract.TypeBrace
	} else if strings.HasPrefix(ln, grammar.TokenLBracket) { // Open bracket.
		l.bracketCount++
		token.Value = grammar.TokenLBracket
		token.Type = fract.TypeBrace
	} else if strings.HasPrefix(ln, grammar.TokenRBracket) { // Close bracket.
		l.bracketCount--
		if l.bracketCount < 0 {
			l.Error("The extra bracket are closed!")
		}
		token.Value = grammar.TokenRBracket
		token.Type = fract.TypeBrace
	} else if strings.HasPrefix(ln, grammar.TokenComma) { // Comma.
		token.Value = grammar.TokenComma
		token.Type = fract.TypeComma
	} else if strings.HasPrefix(ln, grammar.Setter) { // Setter.
		token.Value = grammar.Setter
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.TokenAmper) { // Amper (&).
		token.Value = grammar.TokenAmper
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.TokenVerticalBar) { // Vertical bar (|).
		token.Value = grammar.TokenVerticalBar
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.GreaterEquals) { // Greater than or equals to (>=).
		token.Value = grammar.GreaterEquals
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.LessEquals) { // Less than or equals to (<=).
		token.Value = grammar.LessEquals
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.TokenEquals) { // Equals to (=).
		token.Value = grammar.TokenEquals
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.NotEquals) { // Not equals to (<>).
		token.Value = grammar.NotEquals
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.TokenGreat) { // Greater than (>).
		token.Value = grammar.TokenGreat
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.TokenLess) { // Less than (<).
		token.Value = grammar.TokenLess
		token.Type = fract.TypeOperator
	} else if strings.HasPrefix(ln, grammar.TokenColon) { // Block start.
		token.Value = grammar.TokenColon
		token.Type = fract.TypeBlock
		l.Line--
	} else if isKeywordToken(ln, grammar.KwBlockEnd) { // End of block.
		token.Value = grammar.KwBlockEnd
		token.Type = fract.TypeBlockEnd
	} else if isKeywordToken(ln, grammar.KwVariable) { // Variable.
		token.Value = grammar.KwVariable
		token.Type = fract.TypeVariable
	} else if isKeywordToken(ln, grammar.KwConstVariable) { // Const variable.
		token.Value = grammar.KwConstVariable
		token.Type = fract.TypeVariable
	} else if isKeywordToken(ln, grammar.KwDelete) { // Delete.
		token.Value = grammar.KwDelete
		token.Type = fract.TypeDelete
	} else if isKeywordToken(ln, grammar.KwExit) { // Exit.
		token.Value = grammar.KwExit
		token.Type = fract.TypeExit
	} else if isKeywordToken(ln, grammar.KwIf) { // If.
		token.Value = grammar.KwIf
		token.Type = fract.TypeIf
	} else if isKeywordToken(ln, grammar.KwElseIf) { // Else if.
		token.Value = grammar.KwElseIf
		token.Type = fract.TypeElseIf
	} else if isKeywordToken(ln, grammar.KwForWhileLoop) { // For and while loop.
		token.Value = grammar.KwForWhileLoop
		token.Type = fract.TypeLoop
	} else if isKeywordToken(ln, grammar.KwIn) { // In.
		token.Value = grammar.KwIn
		token.Type = fract.TypeIn
	} else if isKeywordToken(ln, grammar.KwBreak) { // Break.
		token.Value = grammar.KwBreak
		token.Type = fract.TypeBreak
	} else if isKeywordToken(ln, grammar.KwContinue) { // Continue.
		token.Value = grammar.KwContinue
		token.Type = fract.TypeContinue
	} else if isKeywordToken(ln, grammar.KwFunction) { // Function.
		token.Value = grammar.KwFunction
		token.Type = fract.TypeFunction
	} else if isKeywordToken(ln, grammar.DtInt8) { // int8.
		token.Value = grammar.DtInt8
		token.Type = fract.TypeDataType
	} else if isKeywordToken(ln, grammar.DtInt16) { // int16.
		token.Value = grammar.DtInt16
		token.Type = fract.TypeDataType
	} else if isKeywordToken(ln, grammar.DtInt32) { // int32.
		token.Value = grammar.DtInt32
		token.Type = fract.TypeDataType
	} else if isKeywordToken(ln, grammar.DtInt64) { // int64.
		token.Value = grammar.DtInt64
		token.Type = fract.TypeDataType
	} else if isKeywordToken(ln, grammar.DtUInt8) { // uint8.
		token.Value = grammar.DtUInt8
		token.Type = fract.TypeDataType
	} else if isKeywordToken(ln, grammar.DtUInt16) { // uint16.
		token.Value = grammar.DtUInt16
		token.Type = fract.TypeDataType
	} else if isKeywordToken(ln, grammar.DtUInt32) { // uint32.
		token.Value = grammar.DtUInt32
		token.Type = fract.TypeDataType
	} else if isKeywordToken(ln, grammar.DtUInt64) { // uint64.
		token.Value = grammar.DtUInt64
		token.Type = fract.TypeDataType
	} else if isKeywordToken(ln, grammar.DtFloat32) { // float32.
		token.Value = grammar.DtFloat32
		token.Type = fract.TypeDataType
	} else if isKeywordToken(ln, grammar.DtFloat64) { // float64.
		token.Value = grammar.DtFloat64
		token.Type = fract.TypeDataType
	} else if isKeywordToken(ln, grammar.DtBoolean) { // Boolean.
		token.Value = grammar.DtBoolean
		token.Type = fract.TypeDataType
	} else if isKeywordToken(ln, grammar.KwTrue) { // True.
		token.Value = grammar.KwTrue
		token.Type = fract.TypeBooleanTrue
	} else if isKeywordToken(ln, grammar.KwFalse) { // False.
		token.Value = grammar.KwFalse
		token.Type = fract.TypeBooleanFalse
	} else if strings.HasPrefix(ln, grammar.TokenSharp) { // Comment.
	} else if strings.HasPrefix(ln, grammar.TokenQuote) { // Char.
		lexChar(l, &token, fln)
	} else if strings.HasPrefix(ln, grammar.TokenDoubleQuote) { // String.
		lexString(l, &token, fln)
	} else { // Alternates
		/* Check variable name. */
		if check = strings.TrimSpace(regexp.MustCompile(
			"^([A-z])([a-zA-Z0-9" + grammar.TokenUnderscore + grammar.TokenDot +
				"]+)?([[:punct:]]|\\s|$)").FindString(ln)); check != "" { // Name.
			// Remove punct.
			if !strings.HasSuffix(check, grammar.TokenUnderscore) &&
				!strings.HasSuffix(check, grammar.TokenDot) {
				result, _ := regexp.MatchString("(\\s|[[:punct:]])$", check)
				if result {
					check = check[:len(check)-1]
				}
			}

			// Name is finished with dot?
			if strings.HasSuffix(check, grammar.TokenDot) {
				l.Error("What is this?: " + ln)
			}

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
