package lex

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"unicode"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
	"github.com/fract-lang/fract/pkg/str"
)

// Lex of Fract.
type Lex struct {
	lastTk obj.Token

	F            *obj.File
	Col          int // Last column.
	Ln           int // Last line.
	Fin          bool
	RangeComment bool
	Braces       int
	Brackets     int
	Parentheses  int
}

// Error thrown exception.
func (l Lex) Error(msg string) {
	fmt.Printf("File: %s\nPosition: %d:%d\n", l.F.P, l.Ln, l.Col)
	if !l.RangeComment { // Ignore multiline comment error.
		fmt.Println("    " + strings.ReplaceAll(l.F.Lns[l.Ln-1], "\t", " "))
		fmt.Println(str.Whitespace(4+l.Col-2) + "^")
	}
	fmt.Println(msg)
	panic(nil)
}

// Check expected bracket or like and returns true if require retokenize, returns false if not.
// Thrown exception is syntax error.
func (l *Lex) checkExpected(msg string) bool {
	if l.Fin {
		if l.F.P != "<stdin>" {
			l.Ln-- // Subtract for correct line number.
			l.Error(msg)
		}
		return false
	}
	return true
}

// Next lex next line.
func (l *Lex) Next() obj.Tokens {
	var tks obj.Tokens
	// If file is finished?
	if l.Fin {
		return tks
	}
tokenize:
	if l.lastTk.T != fract.StatementTerminator {
		// Restore to defaults.
		l.Col = 1
		l.lastTk.T = fract.None
		l.lastTk.Ln = 0
		l.lastTk.Col = 0
		l.lastTk.V = ""
	}
	// Tokenize line.
	tk := l.Token()
	for tk.T != fract.None {
		if tk.T == fract.StatementTerminator {
			if l.Parentheses == 0 && l.Braces == 0 && l.Brackets == 0 {
				break
			}
			l.Ln++
		}
		if !l.RangeComment && tk.T != fract.Ignore {
			tks = append(tks, tk)
			l.lastTk = tk
		}
		tk = l.Token()
	}
	l.lastTk = tk
	// Go next line.
	l.Ln++
	// Line equals to or bigger then last line.
	l.Fin = l.Ln > len(l.F.Lns)
	switch {
	case l.Parentheses > 0: // Check parentheses.
		if l.checkExpected("Parentheses is expected to close...") {
			goto tokenize
		}
	case l.Braces > 0: // Check braces.
		if l.checkExpected("Brace is expected to close...") {
			goto tokenize
		}
	case l.Brackets > 0: // Check brackets.
		if l.checkExpected("Bracket is expected to close...") {
			goto tokenize
		}
	case l.RangeComment:
		if l.checkExpected("Multiline comment is expected to close...") {
			goto tokenize
		}
	}
	return tks
}

var (
	numericRgx = *regexp.MustCompile(`^(-|)(([0-9]+((\.[0-9]+)|(\.[0-9]+)?(e|E)(\-|\+)[0-9]+)?)|(0x[A-f0-9]+))(\s|[[:punct:]]|$)`)
	macroRgx   = *regexp.MustCompile(`^#(\s+|$)`)
)

// isKeyword returns true if part is keyword, false if not.
func isKeyword(ln, kw string) bool {
	return regexp.MustCompile("^" + kw + `(\s+|$|[[:punct:]])`).MatchString(ln)
}

// isMacro returns true if part is macro, false if not.
func isMacro(ln string) bool { return !macroRgx.MatchString(ln) }

// getName returns name if next token is name, returns empty string if not.
func getName(ln string) string {
	if ln == "" {
		return ln
	}
	for i, r := range ln {
		if r == '-' && i == 0 {
			continue
		} else if r == '.' && i > 0 {
			continue
		} else if r >= '0' && r <= '9' && i > 0 {
			continue
		} else if r == '_' {
			continue
		} else if unicode.IsLetter(r) {
			continue
		}
		if i > 0 {
			return ln[:i]
		}
		return ""
	}
	return ln
}

// getNumeric returns numeric if next token is numeric, returns empty string if not.
func getNumeric(ln string) string { return numericRgx.FindString(ln) }

// Process string espace sequence.
func (l *Lex) strseq(sb *strings.Builder, fln string) bool {
	// Is not espace sequence?
	if fln[l.Col-1] != '\\' {
		return false
	}
	l.Col++
	if l.Col >= len(fln)+1 {
		l.Error("Charray literal is not defined full!")
	}
	switch fln[l.Col-1] {
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

func (l *Lex) lexstr(tk *obj.Token, quote byte, fln string) {
	sb := new(strings.Builder)
	sb.WriteByte(quote)
	l.Col++
	for ; l.Col < len(fln)+1; l.Col++ {
		c := fln[l.Col-1]
		if c == quote { // Finish?
			sb.WriteByte(c)
			break
		} else if !l.strseq(sb, fln) {
			sb.WriteByte(c)
		}
	}
	tk.V = sb.String()
	if tk.V[len(tk.V)-1] != quote {
		l.Error("Close quote is not found!")
	}
	tk.T = fract.Value
	l.Col -= sb.Len() - 1
}

func (l *Lex) lexname(tk *obj.Token, chk string) bool {
	// Remove punct.
	if chk[len(chk)-1] != '_' && chk[len(chk)-1] != '.' {
		r, _ := regexp.MatchString(`(\s|[[:punct:]])$`, chk)
		if r {
			chk = chk[:len(chk)-1]
		}
	}
	// Name is finished with dot?
	if chk[len(chk)-1] == '.' {
		if l.RangeComment { // Ignore comment content.
			l.Col++
			tk.T = fract.Ignore
			return false
		}
		l.Error("What you mean?")
	}
	tk.V = chk
	tk.T = fract.Name
	return true
}

// Generate next token.
func (l *Lex) Token() obj.Token {
	tk := obj.Token{T: fract.None, F: l.F}

	fln := l.F.Lns[l.Ln-1] // Full line.
	// Line is finished.
	if l.Col > len(fln) {
		if l.RangeComment {
			l.F.Lns[l.Ln-1] = ""
		}
		return tk
	}
	// Resume.
	ln := fln[l.Col-1:]
	// Skip spaces.
	for i, c := range ln {
		if unicode.IsSpace(c) {
			l.Col++
			continue
		}
		ln = ln[i:]
		break
	}
	// Content is empty.
	if ln == "" {
		return tk
	}
	// Set token values.
	tk.Col = l.Col
	tk.Ln = l.Ln

	// ************
	//   Tokenize
	// ************

	if l.RangeComment { // Range comment.
		if strings.HasPrefix(ln, "*/") { // Range comment close.
			l.RangeComment = false
			l.Col += 2 // len("<#")
			tk.T = fract.Ignore
			return tk
		}
	}

	switch chk := getNumeric(ln); {
	case (chk != "" &&
		(l.lastTk.V == "" || l.lastTk.T == fract.Operator ||
			(l.lastTk.T == fract.Brace && l.lastTk.V != "]") ||
			l.lastTk.T == fract.StatementTerminator || l.lastTk.T == fract.Loop ||
			l.lastTk.T == fract.Comma || l.lastTk.T == fract.In || l.lastTk.T == fract.If ||
			l.lastTk.T == fract.Else || l.lastTk.T == fract.Ret)) || isKeyword(ln, "NaN"): // Numeric value.
		if chk == "" {
			chk = "NaN"
			l.Col += 3
		} else {
			// Remove punct.
			if lst := chk[len(chk)-1]; lst < '0' || lst > '9' {
				chk = chk[:len(chk)-1]
			}
			l.Col += len(chk)
			if strings.HasPrefix(chk, "0x") {
				// Parse hexadecimal to decimal.
				bi := new(big.Int)
				bi.SetString(chk[2:], 16)
				chk = bi.String()
			} else {
				// Parse floating-point.
				bf := new(big.Float)
				_, f := bf.SetString(chk)
				if !f {
					chk = bf.String()
				}
			}
		}
		tk.V = chk
		tk.T = fract.Value
		return tk
	case strings.HasPrefix(ln, "//"): // Singleline comment.
		l.F.Lns[l.Ln-1] = l.F.Lns[l.Ln-1][:l.Col-1] // Remove comment from original line.
		return tk
	case strings.HasPrefix(ln, "/*"): // Range comment open.
		l.RangeComment = true
		tk.V = "/*"
		tk.T = fract.Ignore
	case ln[0] == '#': // Macro.
		if isMacro(ln) {
			tk.V = "#"
			tk.T = fract.Macro
		}
	case ln[0] == '\'': // String.
		l.lexstr(&tk, '\'', fln)
	case ln[0] == '"': // String.
		l.lexstr(&tk, '"', fln)
	case ln[0] == ';': // Statement terminator.
		tk.V = ";"
		tk.T = fract.StatementTerminator
		l.Ln--
	case strings.HasPrefix(ln, ":="): // Short variable declaration.
		tk.V = ":="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "+="): // Addition assignment.
		tk.V = "+="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "**="): // Exponentiation assignment.
		tk.V = "**="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "*="): // Multiplication assignment.
		tk.V = "*="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "/="): // Division assignment.
		tk.V = "/="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "%="): // Modulus assignment.
		tk.V = "%="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "-="): // Subtraction assignment.
		tk.V = "-="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "<<="): // Left binary shift assignment.
		tk.V = "<<="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, ">>="): // Right binary shift assignment.
		tk.V = ">>="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "|="): // Bitwise Inclusive or assignment.
		tk.V = "|="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "^="): // Bitwise exclusive or assignment.
		tk.V = "^="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "&="): // And assignment.
		tk.V = "&="
		tk.T = fract.Operator
	case ln[0] == '+': // Addition.
		tk.V = "+"
		tk.T = fract.Operator
	case ln[0] == '-': // Subtraction.
		/* Check variable name. */
		if check := getName(ln); check != "" { // Name.
			if !l.lexname(&tk, check) {
				return tk
			}
			break
		}
		tk.V = "-"
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "**"): // Exponentiation.
		tk.V = "**"
		tk.T = fract.Operator
	case ln[0] == '*': // Multiplication.
		tk.V = "*"
		tk.T = fract.Operator
	case ln[0] == '/': // Division.
		tk.V = "/"
		tk.T = fract.Operator
	case ln[0] == '%': // Mod.
		tk.V = "%"
		tk.T = fract.Operator
	case ln[0] == '(': // Open parentheses.
		l.Parentheses++
		tk.V = "("
		tk.T = fract.Brace
	case ln[0] == ')': // Close parentheses.
		l.Parentheses--
		if l.Parentheses < 0 {
			l.Error("The extra parentheses are closed!")
		}
		tk.V = ")"
		tk.T = fract.Brace
	case ln[0] == '{': // Open brace.
		l.Braces++
		tk.V = "{"
		tk.T = fract.Brace
	case ln[0] == '}': // Close brace.
		l.Braces--
		if l.Braces < 0 {
			l.Error("The extra brace are closed!")
		}
		tk.V = "}"
		tk.T = fract.Brace
	case ln[0] == '[': // Open bracket.
		l.Brackets++
		tk.V = "["
		tk.T = fract.Brace
	case ln[0] == ']': // Close bracket.
		l.Brackets--
		if l.Brackets < 0 {
			l.Error("The extra bracket are closed!")
		}
		tk.V = "]"
		tk.T = fract.Brace
	case strings.HasPrefix(ln, "<<"): // Left shift.
		tk.V = "<<"
		tk.T = fract.Operator
	case strings.HasPrefix(ln, ">>"): // Right shift.
		tk.V = ">>"
		tk.T = fract.Operator
	case ln[0] == ',': // Comma.
		tk.V = ","
		tk.T = fract.Comma
	case strings.HasPrefix(ln, "&&"): // Logical and (&&).
		tk.V = "&&"
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "||"): // Logical or (||).
		tk.V = "||"
		tk.T = fract.Operator
	case ln[0] == '|': // Vertical bar.
		tk.V = "|"
		tk.T = fract.Operator
	case ln[0] == '&': // Amper.
		tk.V = "&"
		tk.T = fract.Operator
	case ln[0] == '^': // Bitwise exclusive or(^).
		tk.V = "^"
		tk.T = fract.Operator
	case strings.HasPrefix(ln, ">="): // Greater than or equals to (>=).
		tk.V = ">="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "<="): // Less than or equals to (<=).
		tk.V = "<="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "=="): // Equals to (==).
		tk.V = "=="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "<>"): // Not equals to (<>).
		tk.V = "<>"
		tk.T = fract.Operator
	case ln[0] == '>': // Greater than (>).
		tk.V = ">"
		tk.T = fract.Operator
	case ln[0] == '<': // Less than (<).
		tk.V = "<"
		tk.T = fract.Operator
	case ln[0] == '=': // Equals(=).
		tk.V = "="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "..."): // Params.
		tk.V = "..."
		tk.T = fract.Params
	case isKeyword(ln, "var"): // Variable.
		tk.V = "var"
		tk.T = fract.Var
	case isKeyword(ln, "mut"): // Mutable variable.
		tk.V = "mut"
		tk.T = fract.Var
	case isKeyword(ln, "const"): // Constant.
		tk.V = "const"
		tk.T = fract.Var
	case isKeyword(ln, "protected"): // Protected.
		tk.V = "protected"
		tk.T = fract.Protected
	case isKeyword(ln, "del"): // Delete.
		tk.V = "del"
		tk.T = fract.Delete
	case isKeyword(ln, "defer"): // Defer.
		tk.V = "defer"
		tk.T = fract.Defer
	case isKeyword(ln, "if"): // If.
		tk.V = "if"
		tk.T = fract.If
	case isKeyword(ln, "else"): // Else.
		tk.V = "else"
		tk.T = fract.Else
	case isKeyword(ln, "for"): // Foreach and while loop.
		tk.V = "for"
		tk.T = fract.Loop
	case isKeyword(ln, "in"): // In.
		tk.V = "in"
		tk.T = fract.In
	case isKeyword(ln, "break"): // Break.
		tk.V = "break"
		tk.T = fract.Break
	case isKeyword(ln, "continue"): // Continue.
		tk.V = "continue"
		tk.T = fract.Continue
	case isKeyword(ln, "func"): // Function.
		tk.V = "func"
		tk.T = fract.Func
	case isKeyword(ln, "ret"): // Return.
		tk.V = "ret"
		tk.T = fract.Ret
	case isKeyword(ln, "try"): // Try.
		tk.V = "try"
		tk.T = fract.Try
	case isKeyword(ln, "catch"): // Catch.
		tk.V = "catch"
		tk.T = fract.Catch
	case isKeyword(ln, "open"): // Open.
		tk.V = "open"
		tk.T = fract.Import
	case isKeyword(ln, "true"): // True.
		tk.V = "true"
		tk.T = fract.Value
	case isKeyword(ln, "false"): // False.
		tk.V = "false"
		tk.T = fract.Value
	case isKeyword(ln, "go"): // Concurrency.
		tk.V = "go"
		tk.T = fract.Go
	default: // Alternates
		/* Check variable name. */
		if chk := getName(ln); chk != "" { // Name.
			if !l.lexname(&tk, chk) {
				return tk
			}
		} else { // Error exactly
			if l.RangeComment { // Ignore comment content.
				l.Col++
				tk.T = fract.Ignore
				return tk
			}
			l.Error("What you mean?")
		}
	}

	/* Add length to column. */
	l.Col += len(tk.V)
	return tk
}
