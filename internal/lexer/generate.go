/*
	Generate Function
*/

package lexer

import (
	"regexp"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// processEsacepeSequence Process char literal espace sequence.
// l Lexer.
// sb String builder.
// fln Full line text of current code line.
func processEscapeSequence(l *Lexer, sb *strings.Builder, fln string) bool {
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
		sb.WriteByte('\\')
	case '"':
		sb.WriteByte('"')
	case '\'':
		sb.WriteByte('\'')
	case 'n':
		sb.WriteByte('\n')
	case 'r':
		sb.WriteByte('\r')
	case 't':
		sb.WriteByte('\t')
	case 'b':
		sb.WriteByte('\b')
	case 'f':
		sb.WriteByte('\f')
	case 'a':
		sb.WriteByte('\a')
	case 'v':
		sb.WriteByte('\v')
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
func lexString(l *Lexer, token *obj.Token, quote, fln string) {
	sb := &strings.Builder{}
	sb.WriteString(quote)
	l.Column++
	for ; l.Column < len(fln)+1; l.Column++ {
		char := fln[l.Column-1]
		str := string(char)
		if str == quote { // Finish?
			sb.WriteByte(char)
			break
		} else if !processEscapeSequence(l, sb, fln) {
			sb.WriteByte(char)
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
func (l *Lexer) Generate() obj.Token {
	token := obj.Token{
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

	if l.RangeComment { // Range comment.
		if strings.HasPrefix(ln, grammar.RangeCommentClose) { // Range comment close.
			l.RangeComment = false
			l.Column += len(grammar.RangeCommentClose)
			token.Type = fract.TypeIgnore
			return token
		}
	}

	switch check := strings.TrimSpace(regexp.MustCompile(
		`^(-|)\s*[0-9]+(\.[0-9]+)?(\s|[[:punct:]]|$)`).FindString(ln)); {
	case check != "" &&
		(l.lastToken.Value == "" || l.lastToken.Type == fract.TypeOperator ||
			(l.lastToken.Type == fract.TypeBrace && l.lastToken.Value != grammar.TokenRBracket) ||
			l.lastToken.Type == fract.TypeStatementTerminator || l.lastToken.Type == fract.TypeLoop ||
			l.lastToken.Type == fract.TypeComma || l.lastToken.Type == fract.TypeIn ||
			l.lastToken.Type == fract.TypeIf || l.lastToken.Type == fract.TypeElseIf ||
			l.lastToken.Type == fract.TypeElse || l.lastToken.Type == fract.TypeReturn): // Numeric value.
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
	case strings.HasPrefix(ln, grammar.TokenSemicolon): // Statement terminator.
		token.Value = grammar.TokenSemicolon
		token.Type = fract.TypeStatementTerminator
		l.Line--
	case strings.HasPrefix(ln, grammar.AdditionAssigment): // Addition assigment.
		token.Value = grammar.AdditionAssigment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.MultiplicationAssigment): // Multiplication assigment.
		token.Value = grammar.MultiplicationAssigment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.DivisionAssigment): // Division assigment.
		token.Value = grammar.DivisionAssigment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.PowerAssigment): // Power assigment.
		token.Value = grammar.PowerAssigment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.ModulusAssigment): // Modulus assigment.
		token.Value = grammar.ModulusAssigment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.SubtractionAssigment): // Subtraction assigment.
		token.Value = grammar.SubtractionAssigment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.IntegerDivision): // Integer division.
		token.Value = grammar.IntegerDivision
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.IntegerDivideWithBigger): // Integer divide with bigger.
		token.Value = grammar.IntegerDivideWithBigger
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenPlus): // Addition.
		token.Value = grammar.TokenPlus
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenMinus): // Subtraction.
		token.Value = grammar.TokenMinus
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenStar): // Multiplication.
		token.Value = grammar.TokenStar
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenSlash): // Division.
		token.Value = grammar.TokenSlash
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenCaret): // Exponentiation.
		token.Value = grammar.TokenCaret
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenPercent): // Mod.
		token.Value = grammar.TokenPercent
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenBackslash): // Divisin with bigger.
		token.Value = grammar.TokenBackslash
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenLParenthes): // Open parentheses.
		l.ParenthesCount++
		token.Value = grammar.TokenLParenthes
		token.Type = fract.TypeBrace
	case strings.HasPrefix(ln, grammar.TokenRParenthes): // Close parentheses.
		l.ParenthesCount--
		if l.ParenthesCount < 0 {
			l.Error("The extra parentheses are closed!")
		}
		token.Value = grammar.TokenRParenthes
		token.Type = fract.TypeBrace
	case strings.HasPrefix(ln, grammar.TokenLBrace): // Open brace.
		l.BraceCount++
		token.Value = grammar.TokenLBrace
		token.Type = fract.TypeBrace
	case strings.HasPrefix(ln, grammar.TokenRBrace): // Close brace.
		l.BraceCount--
		if l.BraceCount < 0 {
			l.Error("The extra brace are closed!")
		}
		token.Value = grammar.TokenRBrace
		token.Type = fract.TypeBrace
	case strings.HasPrefix(ln, grammar.TokenLBracket): // Open bracket.
		l.BracketCount++
		token.Value = grammar.TokenLBracket
		token.Type = fract.TypeBrace
	case strings.HasPrefix(ln, grammar.TokenRBracket): // Close bracket.
		l.BracketCount--
		if l.BracketCount < 0 {
			l.Error("The extra bracket are closed!")
		}
		token.Value = grammar.TokenRBracket
		token.Type = fract.TypeBrace
	case strings.HasPrefix(ln, grammar.LeftBinaryShift): // Left shift.
		token.Value = grammar.LeftBinaryShift
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.RightBinaryShift): // Right shift.
		token.Value = grammar.RightBinaryShift
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenComma): // Comma.
		token.Value = grammar.TokenComma
		token.Type = fract.TypeComma
	case strings.HasPrefix(ln, grammar.TokenAmper): // Amper (&).
		token.Value = grammar.TokenAmper
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenVerticalBar): // Vertical bar (|).
		token.Value = grammar.TokenVerticalBar
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.GreaterEquals): // Greater than or equals to (>=).
		token.Value = grammar.GreaterEquals
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.LessEquals): // Less than or equals to (<=).
		token.Value = grammar.LessEquals
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.Equals): // Equals to (==).
		token.Value = grammar.Equals
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.NotEquals): // Not equals to (<>).
		token.Value = grammar.NotEquals
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenGreat): // Greater than (>).
		token.Value = grammar.TokenGreat
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenLess): // Less than (<).
		token.Value = grammar.TokenLess
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenEquals): // Equals(=).
		token.Value = grammar.TokenEquals
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.Params): // Params.
		token.Value = grammar.Params
		token.Type = fract.TypeParams
	case isKeywordToken(ln, grammar.KwBlockEnd): // End of block.
		token.Value = grammar.KwBlockEnd
		token.Type = fract.TypeBlockEnd
	case isKeywordToken(ln, grammar.KwVariable): // Variable.
		token.Value = grammar.KwVariable
		token.Type = fract.TypeVariable
	case isKeywordToken(ln, grammar.KwConstVariable): // Const variable.
		token.Value = grammar.KwConstVariable
		token.Type = fract.TypeVariable
	case isKeywordToken(ln, grammar.KwProtected): // Protected.
		token.Value = grammar.KwProtected
		token.Type = fract.TypeProtected
	case isKeywordToken(ln, grammar.KwDelete): // Delete.
		token.Value = grammar.KwDelete
		token.Type = fract.TypeDelete
	case isKeywordToken(ln, grammar.KwIf): // If.
		token.Value = grammar.KwIf
		token.Type = fract.TypeIf
	case isKeywordToken(ln, grammar.KwElseIf): // Else if.
		token.Value = grammar.KwElseIf
		token.Type = fract.TypeElseIf
	case isKeywordToken(ln, grammar.KwElse): // Else.
		token.Value = grammar.KwElse
		token.Type = fract.TypeElse
	case isKeywordToken(ln, grammar.KwForWhileLoop): // For and while loop.
		token.Value = grammar.KwForWhileLoop
		token.Type = fract.TypeLoop
	case isKeywordToken(ln, grammar.KwIn): // In.
		token.Value = grammar.KwIn
		token.Type = fract.TypeIn
	case isKeywordToken(ln, grammar.KwBreak): // Break.
		token.Value = grammar.KwBreak
		token.Type = fract.TypeBreak
	case isKeywordToken(ln, grammar.KwContinue): // Continue.
		token.Value = grammar.KwContinue
		token.Type = fract.TypeContinue
	case isKeywordToken(ln, grammar.KwFunction): // Function.
		token.Value = grammar.KwFunction
		token.Type = fract.TypeFunction
	case isKeywordToken(ln, grammar.KwReturn): // Return.
		token.Value = grammar.KwReturn
		token.Type = fract.TypeReturn
	case isKeywordToken(ln, grammar.KwTry): // Try.
		token.Value = grammar.KwTry
		token.Type = fract.TypeTry
	case isKeywordToken(ln, grammar.KwCatch): // Catch.
		token.Value = grammar.KwCatch
		token.Type = fract.TypeCatch
	case isKeywordToken(ln, grammar.KwImport): // Open.
		token.Value = grammar.KwImport
		token.Type = fract.TypeImport
	case isKeywordToken(ln, grammar.KwTrue): // True.
		token.Value = grammar.KwTrue
		token.Type = fract.TypeBooleanTrue
	case isKeywordToken(ln, grammar.KwFalse): // False.
		token.Value = grammar.KwFalse
		token.Type = fract.TypeBooleanFalse
	case strings.HasPrefix(ln, grammar.RangeCommentOpen): // Range comment open.
		l.RangeComment = true
		token.Value = grammar.RangeCommentOpen
		token.Type = fract.TypeIgnore
	case strings.HasPrefix(ln, grammar.TokenSharp): // Singleline comment.
		return token
	case strings.HasPrefix(ln, grammar.TokenQuote): // String.
		lexString(l, &token, grammar.TokenQuote, fln)
	case strings.HasPrefix(ln, grammar.TokenDoubleQuote): // String.
		lexString(l, &token, grammar.TokenDoubleQuote, fln)
	default: // Alternates
		/* Check variable name. */
		if check = strings.TrimSpace(regexp.MustCompile(
			`^([A-z])([a-zA-Z0-9_]+)?(\.([a-zA-Z0-9_]+))*([[:punct:]]|\s|$)`).
			FindString(ln)); check != "" { // Name.
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
				if l.RangeComment { // Ignore comment content.
					l.Column++
					token.Type = fract.TypeIgnore
					return token
				}
				l.Error("What you mean?")
			}

			token.Value = check
			token.Type = fract.TypeName
		} else { // Error exactly
			if l.RangeComment { // Ignore comment content.
				l.Column++
				token.Type = fract.TypeIgnore
				return token
			}
			l.Error("What you mean?")
		}
	}

	/* Add length to column. */
	l.Column += len(token.Value)

	return token
}
