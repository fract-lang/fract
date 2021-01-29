package tokenizer

import (
	"fmt"
	"os"
	"strings"

	"../../fract"
	arithmetic "../../fract/arithmetic"
	"../../grammar"
	"../../objects"
)

// Tokenizer Tokenizer of Fract.
type Tokenizer struct {
	/* PUBLIC */

	// Source file.
	File *objects.CodeFile
	// Last column.
	Column int
	// Index of last line.
	Index int
	// Finished destination file.
	Finish bool
}

// New Create new instance of Tokenizer.
// file Destination file.
func New(file *objects.CodeFile) *Tokenizer {
	var _tokenizer *Tokenizer = new(Tokenizer)
	_tokenizer.Column = 1
	_tokenizer.File = file
	return _tokenizer
}

// ExitError Exit with tokenizer error.
// message Message of error.
func (t *Tokenizer) ExitError(message string) {
	fmt.Println()
	fmt.Println("TOKENIZER ERROR")
	fmt.Println("MESSAGE: " + message)
	fmt.Printf("LINE: %d\n", t.File.Lines[t.Index].Line)
	fmt.Printf("COLUMN: %d\n", t.Index)
	os.Exit(1)
}

/* Last putted token */
var lastToken objects.Token

// NextToken Tokenize next token from statement.
func (t *Tokenizer) NextToken() objects.Token {
	var _token objects.Token
	_token.Line = t.Index + 1
	_token.Column = t.Column

	var cline objects.CodeLine = t.File.Lines[t.Index]

	/* Return empty token is statement finished. */
	if t.Column >= len(cline.Text) {
		return _token
	}

	var statement string = cline.Text[t.Column-1:]

	/* Ignore whitespaces and tabs */
	for index := 0; index < len(statement); index++ {
		var ch byte = statement[index]
		if ch == ' ' || ch == 't' {
			t.Column++
		} else {
			statement = statement[index:]
			break
		}
	}

	/* Return empty token if statement is empty. */
	if statement == "" {
		return _token
	}

	/* Arithmetic value check. */
	if statement[0] == grammar.TokenMinus[0] || arithmetic.IsNumeric(statement[0]) {
		var value string = ""
		if statement[0] == grammar.TokenMinus[0] {
			value = grammar.TokenMinus
		}
		if value == "" || (value != "" && (lastToken.Type == fract.TypeOperator ||
			lastToken.Type == fract.TypeOpenParenthes ||
			lastToken.Type == fract.TypeCloseParenthes)) {
			for index := t.Index; index < len(statement); index++ {
				var char byte = statement[index]
				if !arithmetic.IsNumeric(char) && char != grammar.TokenDot[0] {
					break
				}
				value += string(char)
			}
			statement = value
		}
	}

	/* Check anothers. */
	if arithmetic.IsInteger(statement) {
		_token.Type = fract.TypeValue
		_token.Value = statement
	} else if arithmetic.IsFloat(statement) {
		_token.Type = fract.TypeValue
		_token.Value = statement
	} else if strings.HasPrefix(statement, grammar.KwVariable) {
		_token.Type = fract.TypeLet
		_token.Value = grammar.KwVariable
	} else if strings.HasPrefix(statement, grammar.TokenPlus) {
		_token.Type = fract.TypeOperator
		_token.Value = grammar.TokenPlus
	} else if strings.HasPrefix(statement, grammar.TokenMinus) {
		_token.Type = fract.TypeOperator
		_token.Value = grammar.TokenMinus
	} else if strings.HasPrefix(statement, grammar.TokenStar) {
		_token.Type = fract.TypeOperator
		_token.Value = grammar.TokenStar
	} else if strings.HasPrefix(statement, grammar.TokenSlash) {
		_token.Type = fract.TypeOperator
		_token.Value = grammar.TokenSlash
	} else if strings.HasPrefix(statement, grammar.TokenLParenthes) {
		_token.Type = fract.TypeOpenParenthes
		_token.Value = grammar.TokenLParenthes
	} else if strings.HasPrefix(statement, grammar.TokenRParenthes) {
		_token.Type = fract.TypeCloseParenthes
		_token.Value = grammar.TokenRParenthes
	} else {
		t.ExitError("What the?: '" + statement + "'")
	}

	t.Column += len(_token.Value)
	return _token
}

// TokenizeNext Tokenize all statement.
func (t *Tokenizer) TokenizeNext() []objects.Token {
	var tokens []objects.Token

	if t.Finish {
		return tokens
	}

	if t.File.Lines[t.Index].Text == "" {
		return tokens
	}

	/* Reset to defaults */
	t.Column = 1
	lastToken.Type = fract.TypeNone
	lastToken.Value = ""

	var _token objects.Token = t.NextToken()
	for _token.Value != "" {
		tokens = append(tokens, _token)
		lastToken = _token
		_token = t.NextToken()
	}

	if t.Index == len(t.File.Lines)-1 {
		t.Finish = true
	}

	t.Index++
	return tokens
}
