package parser

import (
	"fmt"
	"os"

	"../fract"
	"../fract/arithmetic"
	"../grammar"
	"../objects"
	"../utilities/fs"
	"./formatter"
	"./tokenizer"
)

// Parser Parser of Fract.
type Parser struct {
	/* PRIVATE */

	// Parser of this file.
	file objects.CodeFile
	// Tokenizer of parser.
	tokenizer tokenizer.Tokenizer

	/* PUBLIC */

	/* Type of file. */
	Type int
}

// *********************
//       PRIVATE
// *********************

// printValue Print value to screen.
// value Value to print.
func (p *Parser) printValue(value objects.Value) {
	fmt.Println(value.Content)
}

// processValue Process value from tokens.
// tokens Tokens.
// index Last index.
func (p *Parser) processValue(tokens *[]objects.Token, index *int) objects.Value {
	/* Check parentheses range. */
	for true {
		var result formatter.RangeResult = formatter.LexRange(tokens)
		if result.Found {
			var (
				first  int = 0
				_token objects.Token
			)
			_token.Value = p.processValue(&result.Range, &first).Content
			_token.Type = fract.TypeValue
			*tokens = append(*tokens, *new(objects.Token))
			copy((*tokens)[first+result.Index+1:], (*tokens)[first+result.Index:])
			(*tokens)[first+result.Index] = _token
		} else {
			break
		}
	}

	/*
	* VALUE PROCESS
	 */
	var (
		_value objects.Value
		_type  int = PTypeNone
	)
	for ; *index < len(*tokens); (*index)++ {
		var (
			_token     objects.Token = (*tokens)[*index]
			cacheValue string        = _value.Content
			// cacheType  int           = _value.Type
		)

		/* Check operators. */
		if _token.Value == grammar.TokenPlus {
			_type = PTypeAddition
			continue
		} else if _token.Value == grammar.TokenMinus {
			_type = PTypeSubtraction
			continue
		} else if _token.Value == grammar.TokenStar {
			_type = PTypeMultiplication
			continue
		} else if _token.Value == grammar.TokenSlash {
			_type = PTypeDivision
			continue
		}
		/* End of operator checking */

		/* Check errors */
		if _value.Content != "" && _type == PTypeNone {
			ExitError(_token, "You're write side-by-side two value!")
		}

		/* Value checking */
		if arithmetic.IsInteger(_token.Value) {
			_value.Content = _token.Value
			//_value.type = type_int32;
		} else if arithmetic.IsFloat(_token.Value) {
			_value.Content = _token.Value
			_value.Type = fract.TypeFloat
		} else {
			ExitError(_token, "What the?: "+_token.Value)
		}

		/* If not exists any operator. */
		if _type == PTypeNone {
			continue
		}

		/* If data types are not compatible! */
		/*if(!arithmetic::is_types_compatible(_cache_type, _value.type)) {
		  exit_parser_error(**it, "Data types is not compatible!");
		}*/

		arithmeticValue, err := arithmetic.ToDouble(cacheValue)
		if err != nil {
			ExitError(_token, "Value is not arithmetic!")
		}
		cacheArithmeticValue, err := arithmetic.ToDouble(_value.Content)
		if err != nil {
			ExitError(_token, "Value is not arithmetic!")
		}

		if _type == PTypeAddition {
			_value.Content = arithmetic.FloatToString(arithmeticValue + cacheArithmeticValue)
		} else if _type == PTypeSubtraction {
			_value.Content = arithmetic.FloatToString(arithmeticValue - cacheArithmeticValue)
		} else if _type == PTypeDivision {
			if arithmeticValue == 0 || cacheArithmeticValue == 0 {
				ExitError(_token, "Divide by zero!")
			}
			_value.Content = arithmetic.FloatToString(arithmeticValue / cacheArithmeticValue)
		} else if _type == PTypeMultiplication {
			_value.Content = arithmetic.FloatToString(arithmeticValue * cacheArithmeticValue)
		}

		/* Reset type. */
		_type = PTypeNone
	}

	/* If exists unprocessed operator? */
	if _type != PTypeNone {
		ExitError((*tokens)[(*index)-1], "Unused operator?")
	}

	return _value
}

// checkParentheses Check parentheses.
// tokens Tokens to check.
func (p *Parser) checkParentheses(tokens *[]objects.Token) {
	var (
		count    int = 0
		lastOpen objects.Token
	)

	for index := 0; index < len(*tokens); index++ {
		var _token objects.Token = (*tokens)[index]
		if _token.Type == fract.TypeOpenParenthes {
			lastOpen = _token
			count++
		} else if _token.Type == fract.TypeCloseParenthes {
			if count == 0 {
				ExitError(_token, "The extra parentheses are closed!")
			}
			count--
		}
	}

	if count > 0 {
		ExitError(lastOpen, "The parentheses are opened but not closed!")
	}
}

// *********************
//        PUBLIC
// *********************

// ReadyFile Create instance of code file.
// path Path of file.
func ReadyFile(path string) objects.CodeFile {
	var file objects.CodeFile
	file.Lines = ReadyLines(fs.ReadAllLines(path))
	file.Path = path
	file.File = fs.OpenFile(path)
	return file
}

// ReadyLines Ready lines to process.
// lines Lines to ready.
func ReadyLines(lines []string) []objects.CodeLine {
	var readyLines []objects.CodeLine
	for index := 0; index < len(lines); index++ {
		readyLines = append(readyLines, objects.CodeLine{Line: index + 1, Text: lines[index]})
	}
	return readyLines
}

// ExitError Exit with parser error.
// token Token of error.
// message Message of error.
func ExitError(token objects.Token, message string) {
	fmt.Println()
	fmt.Println("PARSER ERROR")
	fmt.Println("MESSAGE: " + message)
	fmt.Printf("LINE: %d\n", token.Line)
	fmt.Printf("COLUMN: %d\n", token.Column)
	os.Exit(1)
}

// New Create new instance of Parser.
// path Path of destination file.
// type Type of file.
func New(path string, _type int) *Parser {
	var parser *Parser = new(Parser)
	var file objects.CodeFile = ReadyFile(path)
	parser.tokenizer = *tokenizer.New(&file)
	parser.file = *parser.tokenizer.File
	parser.Type = _type
	return parser
}

// Parse Parse code.
func (p *Parser) Parse() {
	for !p.tokenizer.Finish {
		var tokens []objects.Token = p.tokenizer.TokenizeNext()

		p.checkParentheses(&tokens)

		var first objects.Token = tokens[0]
		if first.Type == fract.TypeValue {
			var index int = 0
			p.printValue(p.processValue(&tokens, &index))
		} else {
			ExitError(first, "What the?:"+first.Value)
		}
	}
}
