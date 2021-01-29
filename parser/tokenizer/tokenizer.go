package tokenizer

import (
	"fmt"
	"os"
	"strings"

	arithmeric "../../fract/arithmetic"
	"../../grammar"
	"../../objects"
)

// Tokenizer Tokenizer of Fract.
type Tokenizer struct {
	/* PRIVITE */

	// Source file.
	file *objects.CodeFile
	// Last column.
	column int
	// Index of last line.
	index int

	/* PUBLIC */

	// Finished destination file.
	Finish bool
}

// New Create new instance of Tokenizer.
// file Destination file.
func New(file *objects.CodeFile) *Tokenizer {
	var _tokenizer *Tokenizer = new(Tokenizer)
	_tokenizer.column = 1
	_tokenizer.file = file
	return _tokenizer
}

// ExitError Exit with tokenizer error.
// message Message of error.
func (t *Tokenizer) ExitError(message string) {
	fmt.Println()
	fmt.Println("TOKENIZER ERROR")
	fmt.Println("MESSAGE: " + message)
	fmt.Printf("LINE: %d", t.file.Lines[t.index].Line)
	fmt.Printf("COLUMN: %d", t.index)
	os.Exit(1)
}

/* Last putted token */
var lastToken objects.Token

// NextToken Tokenize next token from statement.
func (t *Tokenizer) NextToken() objects.Token {
	var _token objects.Token
	_token.Line = t.index + 1
	_token.Column = t.column

	var cline objects.CodeLine = t.file.Lines[t.index]

	/* Return empty token is statement finished. */
	if t.column >= len(cline.Text) {
		return _token
	}

	var statement string = cline.Text[t.column-1:]

	/* Ignore whitespaces and tabs */
	for index := 0; index < len(statement); index++ {
		var ch byte = statement[index]
		if ch == ' ' || ch == 't' {
			t.column++
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
	if statement[0] == grammar.TokenMinus[0] || arithmeric.IsNumeric(statement[0]) {
		var value string = ""
		if statement[0] == grammar.TokenMinus[0] {
			value = grammar.TokenMinus
		}
		if value == "" || (
			value != "" && (
				lastToken.Type == TypeOperator      ||
				lastToken.Type == TypeOpenParenthes ||
				lastToken.Type == TypeCloseParenthes
			)
		) {
			for char := range statement[index:] {
				if !arithmeric.IsNumeric(char) && char != grammar.TokenDot[0] {
					break
				}
				value += char
			}
			stastatement = value
		}
	}

	/* Check anothers. */
	if arithmetic.IsInteger(statement) {
		_token.Type = TypeValue
		_token.Value = statement
	} else if arithmetic.IsFloat(statement) {
		_token.Type = TypeValue
		_token.Value = statement
	} else if strings.HasPrefix(statement, grammar.KwVariable) {
		_token.Type = TypeLet
		_token.Value = grammar.KwVariable
	} else if strings.HasPrefix(statement, grammar.TokenPlus) {
		_token.Type = TypeOperator
		_token.Value = grammar.TokenPlus
	} else if strings.HasPrefix(statement, grammar.TokenMinus) {
		_token.Type = TypeOperator
		_token.Value = grammar.TokenMinus
	} else if strings.HasPrefix(statement, grammar.TokenStar) {
		_token.Type = TypeOperator
		_token.Value = grammar.TokenStar
	} else if strings.HasPrefix(statement, grammar.TokenSlash) {
		_token.Type = TypeOperator
		_token.Value = grammar.TokenSlash
	} else if strings.HasPrefix(statement, grammar.TokenLParenthes) {
		_token.Type = TypeOpenParenthes
		_token.Value = grammar.TokenLParenthes
	} else if strings.HasPrefix(statement, grammar.TokenRParenthes) {
		_token.Type = TypeCloseParenthes
		_token.Value = grammar.TokenRParenthes
	} else {
		t.ExitError("What the?: '" + statement + "'")
	}

	t.column += len(_token.Value)
	return _token
}

// TokenizeNext Tokenize all statement.
func (t *Tokenizer) TokenizeNext() []objects.Token {
	var tokens []objects.Token

	if t.Finish {
		return tokens
	}

	if t.file.Lines[t.index].Text == "" {
		return tokens
	}

	/* Reset to defaults */
	t.column = 1
	lastToken.Type = TypeNone
	lastToken.Value = ""

	var _token objects.Token = t.NextToken()
	for _token.Value != "" {
		tokens = append(tokens, _token)
		lastToken = _token
		_token = t.NextToken()
	}

	if t.index == len(t.file.Lines)-1 {
		t.Finish = true
	}

	t.index++
	return tokens
}
