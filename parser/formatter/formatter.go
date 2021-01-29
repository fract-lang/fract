package formatter

import (
	"../../objects"
	"../../utilities/array"
	"../parser"
	"../tokenizer"
)

// LexRange Returns range tokens and remove from original list.
// tokens Tokens to process.
func LexRange(tokens *[]objects.Token) parser.RangeResult {
	var (
		_result parser.RangeResult

		first int
	)

	/* Skip find close parentheses and result ready steps */
	if !_result.found {
		return _result
	}

	/* Find open parentheses */
	for index := 0; index < len(*tokens); index++ {
		var _token objects.Token = (*tokens)[index]
		if _token.Type == tokenizer.TypeOpenParenthes {
			first = index
			_result.Index = index
			_result.Found = true
			break
		}
	}

	/* Find close parentheses */
	var count int = 1
	for index := _result.index + 1; index < len(*tokens); index++ {
		var _token objects.Token = (*tokens)[index]
		if _token.Type == tokenizer.TypeCloseParenthes {
			count = count - 1
			if count == 0 {
				break
			}
		} else if _token.Type == tokenizer.TypeOpenParenthes {
			count = count + 1
		}
		_result.Tokens.PushBack(_token)
	}

	/* Remove range from original tokens */
	array.Remove(tokens, first, len(_result.tokens))

	return _result
}
