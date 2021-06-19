package lexer

// TODO: Add dotted floating point values. Smiliar to: 9. 4. 84.

import (
	"math/big"
	"regexp"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

var (
	numericPattern = *regexp.MustCompile(`^(-|)(([0-9]+((\.[0-9]+)|(\.[0-9]+)?(e|E)(\-|\+)[0-9]+)?)|(0x[A-f0-9]+))(\s|[[:punct:]]|$)`)
	namePattern    = *regexp.MustCompile(`^(-|)([A-z])([a-zA-Z0-9_]+)?(\.([a-zA-Z0-9_]+))*([[:punct:]]|\s|$)`)
	macroPattern   = *regexp.MustCompile(`^#(\s+|$)`)
)

// isKeyword returns true if part is keyword, false if not.
func isKeyword(ln, kw string) bool {
	return regexp.MustCompile("^" + kw + `(\s+|$|[[:punct:]])`).MatchString(ln)
}

// isMacro returns true if part is macro, false if not.
func isMacro(ln string) bool { return !macroPattern.MatchString(ln) }

// getName returns name if next token is name, returns empty string if not.
func getName(ln string) string { return namePattern.FindString(ln) }

// getNumeric returns numeric if next token is numeric, returns empty string if not.
func getNumeric(ln string) string { return numericPattern.FindString(ln) }

// processEsacepeSequence process char literal espace sequence.
func (l *Lexer) processEscapeSequence(sb *strings.Builder, fln string) bool {
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
	if check[len(check)-1] != '_' && check[len(check)-1] != '.' {
		result, _ := regexp.MatchString(`(\s|[[:punct:]])$`, check)
		if result {
			check = check[:len(check)-1]
		}
	}
	// Name is finished with dot?
	if check[len(check)-1] == '.' {
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
	if ln == "" {
		return token
	}
	// Set token values.
	token.Column = l.Column
	token.Line = l.Line

	// ************
	//   Tokenize
	// ************

	if l.RangeComment { // Range comment.
		if strings.HasPrefix(ln, "<#") { // Range comment close.
			l.RangeComment = false
			l.Column += 2 // len("<#")
			token.Type = fract.TypeIgnore
			return token
		}
	}

	switch check := getNumeric(ln); {
	case (check != "" &&
		(l.lastToken.Value == "" || l.lastToken.Type == fract.TypeOperator ||
			(l.lastToken.Type == fract.TypeBrace && l.lastToken.Value != "]") ||
			l.lastToken.Type == fract.TypeStatementTerminator || l.lastToken.Type == fract.TypeLoop ||
			l.lastToken.Type == fract.TypeComma || l.lastToken.Type == fract.TypeIn ||
			l.lastToken.Type == fract.TypeIf || l.lastToken.Type == fract.TypeElseIf ||
			l.lastToken.Type == fract.TypeElse || l.lastToken.Type == fract.TypeReturn)) ||
		isKeyword(ln, "NaN"): // Numeric value.
		if check == "" {
			check = "NaN"
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
	case ln[0] == ';': // Statement terminator.
		token.Value = ";"
		token.Type = fract.TypeStatementTerminator
		l.Line--
	case strings.HasPrefix(ln, "+="): // Addition assignment.
		token.Value = "+="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "**="): // Exponentiation assignment.
		token.Value = "**="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "*="): // Multiplication assignment.
		token.Value = "*="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "/="): // Division assignment.
		token.Value = "/="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "%="): // Modulus assignment.
		token.Value = "%="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "-="): // Subtraction assignment.
		token.Value = "-="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "<<="): // Left binary shift assignment.
		token.Value = "<<="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, ">>="): // Right binary shift assignment.
		token.Value = ">>="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "|="): // Bitwise Inclusive or assignment.
		token.Value = "|="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "^="): // Bitwise exclusive or assignment.
		token.Value = "^="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "&="): // And assignment.
		token.Value = "&="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "//"): // Integer division.
		token.Value = "//"
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "\\\\"): // Integer divide with bigger.
		token.Value = "\\\\"
		token.Type = fract.TypeOperator
	case ln[0] == '+': // Addition.
		token.Value = "+"
		token.Type = fract.TypeOperator
	case ln[0] == '-': // Subtraction.
		/* Check variable name. */
		if check := getName(ln); check != "" { // Name.
			if !l.processName(&token, check) {
				return token
			}
			break
		}
		token.Value = "-"
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "**"): // Exponentiation.
		token.Value = "**"
		token.Type = fract.TypeOperator
	case ln[0] == '*': // Multiplication.
		token.Value = "*"
		token.Type = fract.TypeOperator
	case ln[0] == '/': // Division.
		token.Value = "/"
		token.Type = fract.TypeOperator
	case ln[0] == '%': // Mod.
		token.Value = "%"
		token.Type = fract.TypeOperator
	case ln[0] == '\\': // Divisin with bigger.
		token.Value = "\\"
		token.Type = fract.TypeOperator
	case ln[0] == '(': // Open parentheses.
		l.ParenthesCount++
		token.Value = "("
		token.Type = fract.TypeBrace
	case ln[0] == ')': // Close parentheses.
		l.ParenthesCount--
		if l.ParenthesCount < 0 {
			l.Error("The extra parentheses are closed!")
		}
		token.Value = ")"
		token.Type = fract.TypeBrace
	case ln[0] == '{': // Open brace.
		l.BraceCount++
		token.Value = "{"
		token.Type = fract.TypeBrace
	case ln[0] == '}': // Close brace.
		l.BraceCount--
		if l.BraceCount < 0 {
			l.Error("The extra brace are closed!")
		}
		token.Value = "}"
		token.Type = fract.TypeBrace
	case ln[0] == '[': // Open bracket.
		l.BracketCount++
		token.Value = "["
		token.Type = fract.TypeBrace
	case ln[0] == ']': // Close bracket.
		l.BracketCount--
		if l.BracketCount < 0 {
			l.Error("The extra bracket are closed!")
		}
		token.Value = "]"
		token.Type = fract.TypeBrace
	case strings.HasPrefix(ln, "<<"): // Left shift.
		token.Value = "<<"
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, ">>"): // Right shift.
		token.Value = ">>"
		token.Type = fract.TypeOperator
	case ln[0] == ',': // Comma.
		token.Value = ","
		token.Type = fract.TypeComma
	case strings.HasPrefix(ln, "&&"): // Logical and (&&).
		token.Value = "&&"
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "||"): // Logical or (||).
		token.Value = "||"
		token.Type = fract.TypeOperator
	case ln[0] == '|': // Vertical bar.
		token.Value = "|"
		token.Type = fract.TypeOperator
	case ln[0] == '&': // Amper.
		token.Value = "&"
		token.Type = fract.TypeOperator
	case ln[0] == '^': // Bitwise exclusive or(^).
		token.Value = "^"
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, ">="): // Greater than or equals to (>=).
		token.Value = ">="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "<="): // Less than or equals to (<=).
		token.Value = "<="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "=="): // Equals to (==).
		token.Value = "=="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "<>"): // Not equals to (<>).
		token.Value = "<>"
		token.Type = fract.TypeOperator
	case ln[0] == '>': // Greater than (>).
		token.Value = ">"
		token.Type = fract.TypeOperator
	case ln[0] == '<': // Less than (<).
		token.Value = "<"
		token.Type = fract.TypeOperator
	case ln[0] == '=': // Equals(=).
		token.Value = "="
		token.Type = fract.TypeOperator
	case strings.HasPrefix(ln, "..."): // Params.
		token.Value = "..."
		token.Type = fract.TypeParams
	case isKeyword(ln, "end"): // End of block.
		token.Value = "end"
		token.Type = fract.TypeBlockEnd
	case isKeyword(ln, "var"): // Variable.
		token.Value = "var"
		token.Type = fract.TypeVariable
	case isKeyword(ln, "mut"): // Mutable variable.
		token.Value = "mut"
		token.Type = fract.TypeVariable
	case isKeyword(ln, "const"): // Constant.
		token.Value = "const"
		token.Type = fract.TypeVariable
	case isKeyword(ln, "protected"): // Protected.
		token.Value = "protected"
		token.Type = fract.TypeProtected
	case isKeyword(ln, "del"): // Delete.
		token.Value = "del"
		token.Type = fract.TypeDelete
	case isKeyword(ln, "defer"): // Defer.
		token.Value = "defer"
		token.Type = fract.TypeDefer
	case isKeyword(ln, "if"): // If.
		token.Value = "if"
		token.Type = fract.TypeIf
	case isKeyword(ln, "elif"): // Else if.
		token.Value = "elif"
		token.Type = fract.TypeElseIf
	case isKeyword(ln, "else"): // Else.
		token.Value = "else"
		token.Type = fract.TypeElse
	case isKeyword(ln, "for"): // Foreach and while loop.
		token.Value = "for"
		token.Type = fract.TypeLoop
	case isKeyword(ln, "in"): // In.
		token.Value = "in"
		token.Type = fract.TypeIn
	case isKeyword(ln, "break"): // Break.
		token.Value = "break"
		token.Type = fract.TypeBreak
	case isKeyword(ln, "continue"): // Continue.
		token.Value = "continue"
		token.Type = fract.TypeContinue
	case isKeyword(ln, "func"): // Function.
		token.Value = "func"
		token.Type = fract.TypeFunction
	case isKeyword(ln, "ret"): // Return.
		token.Value = "ret"
		token.Type = fract.TypeReturn
	case isKeyword(ln, "try"): // Try.
		token.Value = "try"
		token.Type = fract.TypeTry
	case isKeyword(ln, "catch"): // Catch.
		token.Value = "catch"
		token.Type = fract.TypeCatch
	case isKeyword(ln, "open"): // Open.
		token.Value = "open"
		token.Type = fract.TypeImport
	case isKeyword(ln, "true"): // True.
		token.Value = "true"
		token.Type = fract.TypeValue
	case isKeyword(ln, "false"): // False.
		token.Value = "false"
		token.Type = fract.TypeValue
	case strings.HasPrefix(ln, "#>"): // Range comment open.
		l.RangeComment = true
		token.Value = "#>"
		token.Type = fract.TypeIgnore
	case ln[0] == '#': // Singleline comment or macro.
		if isMacro(ln) {
			token.Value = "#"
			token.Type = fract.TypeMacro
		} else {
			l.File.Lines[l.Line-1] = l.File.Lines[l.Line-1][:l.Column-1] // Remove comment from original line.
			return token
		}
	case ln[0] == '\'': // String.
		l.lexString(&token, '\'', fln)
	case ln[0] == '"': // String.
		l.lexString(&token, '"', fln)
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
