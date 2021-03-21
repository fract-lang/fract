/*
	Generate Function
*/

package lexer

import (
	"regexp"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
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

// lexString Lex string literal.
// l Lexer.
// token Token.
// quote Quote style.
// fln Full line text of current code line.
func lexString(l *Lexer, token *objects.Token, quote, fln string) {
	var sb strings.Builder
	sb.WriteString(quote)
	l.Column++
	for ; l.Column < len(fln)+1; l.Column++ {
		current := string(fln[l.Column-1])
		if current == quote { // Finish?
			sb.WriteString(current)
			break
		} else if !processEscapeSequence(l, token, fln) {
			sb.WriteString(current)
		}
	}
	token.Value = sb.String()
	if !strings.HasSuffix(token.Value, quote) {
		l.Error("Close quote is not found!")
	}
	token.Type = fract.TypeValue

	l.Column -= sb.Len() - 1
}

// Generate Generate next token.
func (l *Lexer) Generate() objects.Token {
	token := objects.Token{
		Type: fract.TypeNone,
		File: l.File,
	}
	fln := l.File.Lines[l.Line-1].Text // Full line.

	/* Line is finished. */
	if l.Column > len(fln) {
		return token
	}

	// Resume.
	ln := fln[l.Column-1:]

	/* Skip spaces. */
	for index, char := range ln {
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

	if l.multilineComment { // Multiline comment.
		if strings.HasPrefix(ln, grammar.MultilineCommentClose) { // Multiline comment close.
			l.multilineComment = false
			l.Column += len(grammar.MultilineCommentClose)
			token.Type = fract.TypeIgnore
			return token
		}
	}

	/* Check arithmetic value? */
	if check := strings.TrimSpace(regexp.MustCompile(
		`^(-|)\s*[0-9]+(\.[0-9]+)?(\s|[[:punct:]]|$)`).FindString(ln)); check != "" &&
		(l.lastToken.Value == "" || l.lastToken.Type == fract.TypeOperator ||
			(l.lastToken.Type == fract.TypeBrace && l.lastToken.Value != grammar.TokenRBracket) ||
			l.lastToken.Type == fract.TypeStatementTerminator || l.lastToken.Type == fract.TypeLoop ||
			l.lastToken.Type == fract.TypeComma || l.lastToken.Type == fract.TypeIn ||
			l.lastToken.Type == fract.TypeIf || l.lastToken.Type == fract.TypeElseIf ||
			l.lastToken.Type == fract.TypeElse || l.lastToken.Type == fract.TypeExit ||
			l.lastToken.Type == fract.TypeReturn) { // Numeric value.
		// Remove punct.
		result, _ := regexp.MatchString(`(\s|[[:punct:]])$`, check)
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
	} else if strings.HasPrefix(ln, grammar.Input) { // Input (<<).
		token.Value = grammar.Input
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
	} else if strings.HasPrefix(ln, grammar.Equals) { // Equals to (==).
		token.Value = grammar.Equals
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
	} else if strings.HasPrefix(ln, grammar.TokenEquals) { // Equals(=).
		token.Value = grammar.TokenEquals
		token.Type = fract.TypeOperator
	} else if isKeywordToken(ln, grammar.KwBlockEnd) { // End of block.
		token.Value = grammar.KwBlockEnd
		token.Type = fract.TypeBlockEnd
	} else if isKeywordToken(ln, grammar.KwVariable) { // Variable.
		token.Value = grammar.KwVariable
		token.Type = fract.TypeVariable
	} else if isKeywordToken(ln, grammar.KwConstVariable) { // Const variable.
		token.Value = grammar.KwConstVariable
		token.Type = fract.TypeVariable
	} else if isKeywordToken(ln, grammar.KwProtected) { // Protected.
		token.Value = grammar.KwProtected
		token.Type = fract.TypeProtected
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
	} else if isKeywordToken(ln, grammar.KwElse) { // Else.
		token.Value = grammar.KwElse
		token.Type = fract.TypeElse
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
	} else if isKeywordToken(ln, grammar.KwReturn) { // Return.
		token.Value = grammar.KwReturn
		token.Type = fract.TypeReturn
	} else if isKeywordToken(ln, grammar.KwTrue) { // True.
		token.Value = grammar.KwTrue
		token.Type = fract.TypeBooleanTrue
	} else if isKeywordToken(ln, grammar.KwFalse) { // False.
		token.Value = grammar.KwFalse
		token.Type = fract.TypeBooleanFalse
	} else if strings.HasPrefix(ln, grammar.MultilineCommentOpen) { // Multiline comment open.
		l.multilineComment = true
		token.Value = grammar.MultilineCommentOpen
		token.Type = fract.TypeIgnore
	} else if strings.HasPrefix(ln, grammar.TokenSharp) { // Singleline comment.
		return token
	} else if strings.HasPrefix(ln, grammar.TokenQuote) { // String.
		lexString(l, &token, grammar.TokenQuote, fln)
	} else if strings.HasPrefix(ln, grammar.TokenDoubleQuote) { // String.
		lexString(l, &token, grammar.TokenDoubleQuote, fln)
	} else { // Alternates
		/* Check variable name. */
		if check = strings.TrimSpace(regexp.MustCompile(
			`^([A-z])([a-zA-Z0-9` + grammar.TokenUnderscore + grammar.TokenDot +
				`]+)?([[:punct:]]|\s|$)`).FindString(ln)); check != "" { // Name.
			// Remove punct.
			if !strings.HasSuffix(check, grammar.TokenUnderscore) &&
				!strings.HasSuffix(check, grammar.TokenDot) {
				result, _ := regexp.MatchString(`(\s|[[:punct:]])$`, check)
				if result {
					check = check[:len(check)-1]
				}
			}

			// Name is finished with dot?
			if strings.HasSuffix(check, grammar.TokenDot) {
				if l.multilineComment { // Ignore comment content.
					l.Column++
					token.Type = fract.TypeIgnore
					return token
				}
				l.Error("What is this?")
			}

			token.Value = strings.TrimSpace(check)
			token.Type = fract.TypeName
		} else { // Error exactly
			if l.multilineComment { // Ignore comment content.
				l.Column++
				token.Type = fract.TypeIgnore
				return token
			}
			l.Error("What is this?")
		}
	}

	/* Add length to column. */
	l.Column += len(token.Value)

	return token
}
