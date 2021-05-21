package lexer

import (
	"math/big"
	"regexp"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
)

var (
	numericPattern = *regexp.MustCompile(`^(-|)(([0-9]+((\.[0-9]+)|(\.[0-9]+)?(e|E)(\-|\+)[0-9]+)?)|(0x[A-f0-9]+))(\s|[[:punct:]]|$)`)
	namePattern    = *regexp.MustCompile(`^(-|)([A-z])([a-zA-Z0-9_]+)?(\.([a-zA-Z0-9_]+))*([[:punct:]]|\s|$)`)
	macroPattern   = *regexp.MustCompile(`^#(\s+|$)`)
)

// isKeyword returns true if part is keyword, false if not.
func isKeyword(ln, kw string) bool { return regexp.MustCompile("^" + kw  + `(\s+|$|[[:punct:]])`).MatchString(ln) }

// isMacro returns true if part is macro, false if not.
func isMacro(ln string) bool { return !macroPattern.MatchString(ln) }

// getName returns name if next token is name, returns empty string if not.
func getName(ln string) string { return namePattern.FindString(ln) }

// getNumeric returns numeric if next token is numeric, returns empty string if not.
func getNumeric(ln string) string { return numericPattern.FindString(ln) }

// processEsacepeSequence process char literal espace sequence.
func (l *Lexer) processEscapeSequence(sb *strings.Builder, fln string) bool {
	// Is not espace sequence?
	if fln[l.Column-1] != '\\' { return false }

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

// lexString lex string literal.
func (l *Lexer) lexString(token *objects.Token, quote byte, fln string) {
	sb := new(strings.Builder)
	sb.WriteByte(quote)
	l.Column++
	for ; l.Column < len(fln)+1; l.Column++ {
		char := fln[l.Column-1]
		if char == quote { // Finish?
			sb.WriteByte(char)
			break
		} else if !l.processEscapeSequence(sb, fln) {
			sb.WriteByte(char)
		}
	}
	token.Value = sb.String()
	if token.Value[len(token.Value)-1] != quote {
		l.Error("Close quote is not found!")
	}
	token.Type = fract.TypeValue

	l.Column -= sb.Len() - 1
}

func (l *Lexer) processName(token *objects.Token, check string) bool {
	// Remove punct.
	if !strings.HasSuffix(check, grammar.TokenUnderscore) && !strings.HasSuffix(check, grammar.TokenDot) {
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
			return false
		}
		l.Error("What you mean?")
	}

	token.Value = check
	token.Type = fract.TypeName
	return true
}

// Generate next token.
func (l *Lexer) Generate() objects.Token {
	token := objects.Token{
		Type: fract.TypeNone,
		File: l.File,
	}

	fln := l.File.Lines[l.Line-1] // Full line.

	// Line is finished.
	if l.Column > len(fln) {
		if l.RangeComment {
			l.File.Lines[l.Line-1] = ""
		}
		return token
	}

	// Resume.
	ln := fln[l.Column-1:]

	// Skip spaces.
	for index, char := range ln {
		if char == ' ' || char == '\t' {
			l.Column++
			continue
		}
		ln = ln[index:]
		break
	}

	// Content is empty.
	if ln == "" { return token }

	// Set token values.
	token.Column = l.Column
	token.Line = l.Line

	// ************
	//   Tokenize
	// ************

	if l.RangeComment { // Range comment.
		if strings.HasPrefix(ln, grammar.RangeCommentClose) { // Range comment close.
			l.RangeComment = false
			l.Column += len(grammar.RangeCommentClose)
			token.Type = fract.TypeIgnore
			return token
		}
	}

	switch check := getNumeric(ln); {
	case (check != "" &&
		(l.lastToken.Value == "" || l.lastToken.Type == fract.TypeOperator ||
			(l.lastToken.Type == fract.TypeBrace && l.lastToken.Value != grammar.TokenRBracket) ||
			l.lastToken.Type == fract.TypeStatementTerminator || l.lastToken.Type == fract.TypeLoop ||
			l.lastToken.Type == fract.TypeComma || l.lastToken.Type == fract.TypeIn ||
			l.lastToken.Type == fract.TypeIf || l.lastToken.Type == fract.TypeElseIf ||
			l.lastToken.Type == fract.TypeElse || l.lastToken.Type == fract.TypeReturn)) ||
		isKeyword(ln, grammar.KwNaN): // Numeric value.
		if check == "" {
			check = grammar.KwNaN
			l.Column += 3
		} else {
			// Remove punct.
			if last := check[len(check)-1]; last != '0' && last != '1' &&
				last != '2' && last != '3' && last != '4' && last != '5' &&
				last != '6' && last != '7' && last != '8' && last != '9' {
				check = check[:len(check)-1]
			}

			l.Column += len(check)

			if strings.HasPrefix(check, "0x") {
				// Parse hexadecimal to decimal.
				bigInt := new(big.Int)
				bigInt.SetString(check[2:], 16)
				check = bigInt.String()
			} else {
				// Parse floating-point.
				bigFloat := new(big.Float)
				_, fail := bigFloat.SetString(check)
				if !fail {
					check = bigFloat.String()
				}
			}
		}

		token.Value = check
		token.Type = fract.TypeValue
		return token
	case strings.HasPrefix(ln, grammar.TokenSemicolon): // Statement terminator.
		token.Value = grammar.TokenSemicolon
		token.Type = fract.TypeStatementTerminator
		l.Line--
	case strings.HasPrefix(ln, grammar.AdditionAssignment): // Addition assignment.
		token.Value = grammar.AdditionAssignment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.ExponentiationAssignment): // Exponentiation assignment.
		token.Value = grammar.ExponentiationAssignment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.MultiplicationAssignment): // Multiplication assignment.
		token.Value = grammar.MultiplicationAssignment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.DivisionAssignment): // Division assignment.
		token.Value = grammar.DivisionAssignment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.ModulusAssignment): // Modulus assignment.
		token.Value = grammar.ModulusAssignment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.SubtractionAssignment): // Subtraction assignment.
		token.Value = grammar.SubtractionAssignment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.LeftBinaryShiftAssignment): // Left binary shift assignment.
		token.Value = grammar.LeftBinaryShiftAssignment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.RightBinaryShiftAssignment): // Right binary shift assignment.
		token.Value = grammar.RightBinaryShiftAssignment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.InclusiveOrAssignment): // Bitwise Inclusive or assignment.
		token.Value = grammar.InclusiveOrAssignment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.XOrAssignment): // Bitwise exclusive or assignment.
		token.Value = grammar.XOrAssignment
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.AndAssignment): // And assignment.
		token.Value = grammar.AndAssignment
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
		/* Check variable name. */
		if check := getName(ln); check != "" { // Name.
			if !l.processName(&token, check) {
				return token
			}
			break
		}
		token.Value = grammar.TokenMinus
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.Exponentiation): // Exponentiation.
		token.Value = grammar.Exponentiation
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenStar): // Multiplication.
		token.Value = grammar.TokenStar
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenSlash): // Division.
		token.Value = grammar.TokenSlash
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
	case strings.HasPrefix(ln, grammar.LogicalAnd): // Logical and (&&).
		token.Value = grammar.LogicalAnd
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.LogicalOr): // Logical or (||).
		token.Value = grammar.LogicalOr
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenVerticalBar): // Vertical bar.
		token.Value = grammar.TokenVerticalBar
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenAmper): // Amper.
		token.Value = grammar.TokenAmper
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, grammar.TokenCaret): // Bitwise exclusive or(^).
		token.Value = grammar.TokenCaret
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
	case isKeyword(ln, grammar.KwBlockEnd): // End of block.
		token.Value = grammar.KwBlockEnd
		token.Type = fract.TypeBlockEnd
	case isKeyword(ln, grammar.KwVariable): // Variable.
		token.Value = grammar.KwVariable
		token.Type = fract.TypeVariable
	case isKeyword(ln, grammar.KwConstVariable): // Const variable.
		token.Value = grammar.KwConstVariable
		token.Type = fract.TypeVariable
	case isKeyword(ln, grammar.KwProtected): // Protected.
		token.Value = grammar.KwProtected
		token.Type = fract.TypeProtected
	case isKeyword(ln, grammar.KwDelete): // Delete.
		token.Value = grammar.KwDelete
		token.Type = fract.TypeDelete
	case isKeyword(ln, grammar.KwIf): // If.
		token.Value = grammar.KwIf
		token.Type = fract.TypeIf
	case isKeyword(ln, grammar.KwElseIf): // Else if.
		token.Value = grammar.KwElseIf
		token.Type = fract.TypeElseIf
	case isKeyword(ln, grammar.KwElse): // Else.
		token.Value = grammar.KwElse
		token.Type = fract.TypeElse
	case isKeyword(ln, grammar.KwForWhileLoop): // For and while loop.
		token.Value = grammar.KwForWhileLoop
		token.Type = fract.TypeLoop
	case isKeyword(ln, grammar.KwIn): // In.
		token.Value = grammar.KwIn
		token.Type = fract.TypeIn
	case isKeyword(ln, grammar.KwBreak): // Break.
		token.Value = grammar.KwBreak
		token.Type = fract.TypeBreak
	case isKeyword(ln, grammar.KwContinue): // Continue.
		token.Value = grammar.KwContinue
		token.Type = fract.TypeContinue
	case isKeyword(ln, grammar.KwFunction): // Function.
		token.Value = grammar.KwFunction
		token.Type = fract.TypeFunction
	case isKeyword(ln, grammar.KwReturn): // Return.
		token.Value = grammar.KwReturn
		token.Type = fract.TypeReturn
	case isKeyword(ln, grammar.KwTry): // Try.
		token.Value = grammar.KwTry
		token.Type = fract.TypeTry
	case isKeyword(ln, grammar.KwCatch): // Catch.
		token.Value = grammar.KwCatch
		token.Type = fract.TypeCatch
	case isKeyword(ln, grammar.KwImport): // Open.
		token.Value = grammar.KwImport
		token.Type = fract.TypeImport
	case isKeyword(ln, grammar.KwTrue): // True.
		token.Value = grammar.KwTrue
		token.Type = fract.TypeValue
	case isKeyword(ln, grammar.KwFalse): // False.
		token.Value = grammar.KwFalse
		token.Type = fract.TypeValue
	case strings.HasPrefix(ln, grammar.RangeCommentOpen): // Range comment open.
		l.RangeComment = true
		token.Value = grammar.RangeCommentOpen
		token.Type = fract.TypeIgnore
	case strings.HasPrefix(ln, grammar.TokenSharp): // Singleline comment or macro.
		if isMacro(ln) {
			token.Value = grammar.TokenSharp
			token.Type = fract.TypeMacro
		} else {
			l.File.Lines[l.Line-1] = l.File.Lines[l.Line-1][:l.Column-1] // Remove comment from original line.
			return token
		}
	case strings.HasPrefix(ln, grammar.TokenQuote): // String.
		l.lexString(&token, grammar.TokenQuote[0], fln)
	case strings.HasPrefix(ln, grammar.TokenDoubleQuote): // String.
		l.lexString(&token, grammar.TokenDoubleQuote[0], fln)
	default: // Alternates
		/* Check variable name. */
		if check := getName(ln); check != "" { // Name.
			if !l.processName(&token, check) {
				return token
			}
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
