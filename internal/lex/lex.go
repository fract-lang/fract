package lex

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

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
func (l *Lex) Next() []obj.Token {
	var tks []obj.Token
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
		l.lastTk.Val = ""
	}
	// Tokenize line.
	tk := l.Generate()
	for tk.T != fract.None &&
		tk.T != fract.StatementTerminator {
		if !l.RangeComment && tk.T != fract.Ignore {
			tks = append(tks, tk)
			l.lastTk = tk
		}
		tk = l.Generate()
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
	numericRegexp = *regexp.MustCompile(`^(-|)(([0-9]+((\.[0-9]+)|(\.[0-9]+)?(e|E)(\-|\+)[0-9]+)?)|(0x[A-f0-9]+))(\s|[[:punct:]]|$)`)
	nameRegexp    = *regexp.MustCompile(`^(-|)([A-z])([a-zA-Z0-9_]+)?(\.([a-zA-Z0-9_]+))*([[:punct:]]|\s|$)`)
	macroRegexp   = *regexp.MustCompile(`^#(\s+|$)`)
)

// isKeyword returns true if part is keyword, false if not.
func isKeyword(ln, kw string) bool {
	return regexp.MustCompile("^" + kw + `(\s+|$|[[:punct:]])`).MatchString(ln)
}

// isMacro returns true if part is macro, false if not.
func isMacro(ln string) bool { return !macroRegexp.MatchString(ln) }

// getName returns name if next token is name, returns empty string if not.
func getName(ln string) string { return nameRegexp.FindString(ln) }

// getNumeric returns numeric if next token is numeric, returns empty string if not.
func getNumeric(ln string) string { return numericRegexp.FindString(ln) }

// processEsacepeSequence process char literal espace sequence.
func (l *Lex) processEscapeSequence(sb *strings.Builder, fln string) bool {
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

// lexString lex string literal.
func (l *Lex) lexString(tk *obj.Token, quote byte, fln string) {
	sb := new(strings.Builder)
	sb.WriteByte(quote)
	l.Col++
	for ; l.Col < len(fln)+1; l.Col++ {
		c := fln[l.Col-1]
		if c == quote { // Finish?
			sb.WriteByte(c)
			break
		} else if !l.processEscapeSequence(sb, fln) {
			sb.WriteByte(c)
		}
	}
	tk.Val = sb.String()
	if tk.Val[len(tk.Val)-1] != quote {
		l.Error("Close quote is not found!")
	}
	tk.T = fract.Value
	l.Col -= sb.Len() - 1
}

func (l *Lex) processName(tk *obj.Token, chk string) bool {
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
	tk.Val = chk
	tk.T = fract.Name
	return true
}

// Generate next token.
func (l *Lex) Generate() obj.Token {
	tk := obj.Token{
		T: fract.None,
		F: l.F,
	}

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
		if c == ' ' || c == '\t' {
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
		if strings.HasPrefix(ln, "<#") { // Range comment close.
			l.RangeComment = false
			l.Col += 2 // len("<#")
			tk.T = fract.Ignore
			return tk
		}
	}

	switch chk := getNumeric(ln); {
	case (chk != "" &&
		(l.lastTk.Val == "" || l.lastTk.T == fract.Operator ||
			(l.lastTk.T == fract.Brace && l.lastTk.Val != "]") ||
			l.lastTk.T == fract.StatementTerminator || l.lastTk.T == fract.Loop ||
			l.lastTk.T == fract.Comma || l.lastTk.T == fract.In ||
			l.lastTk.T == fract.If || l.lastTk.T == fract.ElseIf ||
			l.lastTk.T == fract.Else || l.lastTk.T == fract.Ret)) ||
		isKeyword(ln, "NaN"): // Numeric value.
		if chk == "" {
			chk = "NaN"
			l.Col += 3
		} else {
			// Remove punct.
			if lst := chk[len(chk)-1]; lst != '0' && lst != '1' &&
				lst != '2' && lst != '3' && lst != '4' && lst != '5' &&
				lst != '6' && lst != '7' && lst != '8' && lst != '9' {
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
		tk.Val = chk
		tk.T = fract.Value
		return tk
	case ln[0] == ';': // Statement terminator.
		tk.Val = ";"
		tk.T = fract.StatementTerminator
		l.Ln--
	case strings.HasPrefix(ln, "+="): // Addition assignment.
		tk.Val = "+="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "**="): // Exponentiation assignment.
		tk.Val = "**="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "*="): // Multiplication assignment.
		tk.Val = "*="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "/="): // Division assignment.
		tk.Val = "/="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "%="): // Modulus assignment.
		tk.Val = "%="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "-="): // Subtraction assignment.
		tk.Val = "-="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "<<="): // Left binary shift assignment.
		tk.Val = "<<="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, ">>="): // Right binary shift assignment.
		tk.Val = ">>="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "|="): // Bitwise Inclusive or assignment.
		tk.Val = "|="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "^="): // Bitwise exclusive or assignment.
		tk.Val = "^="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "&="): // And assignment.
		tk.Val = "&="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "//"): // Integer division.
		tk.Val = "//"
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "\\\\"): // Integer divide with bigger.
		tk.Val = "\\\\"
		tk.T = fract.Operator
	case ln[0] == '+': // Addition.
		tk.Val = "+"
		tk.T = fract.Operator
	case ln[0] == '-': // Subtraction.
		/* Check variable name. */
		if check := getName(ln); check != "" { // Name.
			if !l.processName(&tk, check) {
				return tk
			}
			break
		}
		tk.Val = "-"
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "**"): // Exponentiation.
		tk.Val = "**"
		tk.T = fract.Operator
	case ln[0] == '*': // Multiplication.
		tk.Val = "*"
		tk.T = fract.Operator
	case ln[0] == '/': // Division.
		tk.Val = "/"
		tk.T = fract.Operator
	case ln[0] == '%': // Mod.
		tk.Val = "%"
		tk.T = fract.Operator
	case ln[0] == '\\': // Divisin with bigger.
		tk.Val = "\\"
		tk.T = fract.Operator
	case ln[0] == '(': // Open parentheses.
		l.Parentheses++
		tk.Val = "("
		tk.T = fract.Brace
	case ln[0] == ')': // Close parentheses.
		l.Parentheses--
		if l.Parentheses < 0 {
			l.Error("The extra parentheses are closed!")
		}
		tk.Val = ")"
		tk.T = fract.Brace
	case ln[0] == '{': // Open brace.
		l.Braces++
		tk.Val = "{"
		tk.T = fract.Brace
	case ln[0] == '}': // Close brace.
		l.Braces--
		if l.Braces < 0 {
			l.Error("The extra brace are closed!")
		}
		tk.Val = "}"
		tk.T = fract.Brace
	case ln[0] == '[': // Open bracket.
		l.Brackets++
		tk.Val = "["
		tk.T = fract.Brace
	case ln[0] == ']': // Close bracket.
		l.Brackets--
		if l.Brackets < 0 {
			l.Error("The extra bracket are closed!")
		}
		tk.Val = "]"
		tk.T = fract.Brace
	case strings.HasPrefix(ln, "<<"): // Left shift.
		tk.Val = "<<"
		tk.T = fract.Operator
	case strings.HasPrefix(ln, ">>"): // Right shift.
		tk.Val = ">>"
		tk.T = fract.Operator
	case ln[0] == ',': // Comma.
		tk.Val = ","
		tk.T = fract.Comma
	case strings.HasPrefix(ln, "&&"): // Logical and (&&).
		tk.Val = "&&"
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "||"): // Logical or (||).
		tk.Val = "||"
		tk.T = fract.Operator
	case ln[0] == '|': // Vertical bar.
		tk.Val = "|"
		tk.T = fract.Operator
	case ln[0] == '&': // Amper.
		tk.Val = "&"
		tk.T = fract.Operator
	case ln[0] == '^': // Bitwise exclusive or(^).
		tk.Val = "^"
		tk.T = fract.Operator
	case strings.HasPrefix(ln, ">="): // Greater than or equals to (>=).
		tk.Val = ">="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "<="): // Less than or equals to (<=).
		tk.Val = "<="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "=="): // Equals to (==).
		tk.Val = "=="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "<>"): // Not equals to (<>).
		tk.Val = "<>"
		tk.T = fract.Operator
	case ln[0] == '>': // Greater than (>).
		tk.Val = ">"
		tk.T = fract.Operator
	case ln[0] == '<': // Less than (<).
		tk.Val = "<"
		tk.T = fract.Operator
	case ln[0] == '=': // Equals(=).
		tk.Val = "="
		tk.T = fract.Operator
	case strings.HasPrefix(ln, "..."): // Params.
		tk.Val = "..."
		tk.T = fract.Params
	case isKeyword(ln, "end"): // End of block.
		tk.Val = "end"
		tk.T = fract.End
	case isKeyword(ln, "var"): // Variable.
		tk.Val = "var"
		tk.T = fract.Var
	case isKeyword(ln, "mut"): // Mutable variable.
		tk.Val = "mut"
		tk.T = fract.Var
	case isKeyword(ln, "const"): // Constant.
		tk.Val = "const"
		tk.T = fract.Var
	case isKeyword(ln, "protected"): // Protected.
		tk.Val = "protected"
		tk.T = fract.Protected
	case isKeyword(ln, "del"): // Delete.
		tk.Val = "del"
		tk.T = fract.Delete
	case isKeyword(ln, "defer"): // Defer.
		tk.Val = "defer"
		tk.T = fract.Defer
	case isKeyword(ln, "if"): // If.
		tk.Val = "if"
		tk.T = fract.If
	case isKeyword(ln, "elif"): // Else if.
		tk.Val = "elif"
		tk.T = fract.ElseIf
	case isKeyword(ln, "else"): // Else.
		tk.Val = "else"
		tk.T = fract.Else
	case isKeyword(ln, "for"): // Foreach and while loop.
		tk.Val = "for"
		tk.T = fract.Loop
	case isKeyword(ln, "in"): // In.
		tk.Val = "in"
		tk.T = fract.In
	case isKeyword(ln, "break"): // Break.
		tk.Val = "break"
		tk.T = fract.Break
	case isKeyword(ln, "continue"): // Continue.
		tk.Val = "continue"
		tk.T = fract.Continue
	case isKeyword(ln, "func"): // Function.
		tk.Val = "func"
		tk.T = fract.Func
	case isKeyword(ln, "ret"): // Return.
		tk.Val = "ret"
		tk.T = fract.Ret
	case isKeyword(ln, "try"): // Try.
		tk.Val = "try"
		tk.T = fract.Try
	case isKeyword(ln, "catch"): // Catch.
		tk.Val = "catch"
		tk.T = fract.Catch
	case isKeyword(ln, "open"): // Open.
		tk.Val = "open"
		tk.T = fract.Import
	case isKeyword(ln, "true"): // True.
		tk.Val = "true"
		tk.T = fract.Value
	case isKeyword(ln, "false"): // False.
		tk.Val = "false"
		tk.T = fract.Value
	case strings.HasPrefix(ln, "#>"): // Range comment open.
		l.RangeComment = true
		tk.Val = "#>"
		tk.T = fract.Ignore
	case ln[0] == '#': // Singleline comment or macro.
		if isMacro(ln) {
			tk.Val = "#"
			tk.T = fract.Macro
		} else {
			l.F.Lns[l.Ln-1] = l.F.Lns[l.Ln-1][:l.Col-1] // Remove comment from original line.
			return tk
		}
	case ln[0] == '\'': // String.
		l.lexString(&tk, '\'', fln)
	case ln[0] == '"': // String.
		l.lexString(&tk, '"', fln)
	default: // Alternates
		/* Check variable name. */
		if chk := getName(ln); chk != "" { // Name.
			if !l.processName(&tk, chk) {
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
	l.Col += len(tk.Val)
	return tk
}
