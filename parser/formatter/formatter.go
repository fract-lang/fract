package formatter

import (
	"../../fract"
	"../../objects"
)

// LexRange Returns range tokens and remove from original list.
// tokens Tokens to process.
func LexRange(tokens *[]objects.Token) RangeResult {
	var (
		_result RangeResult

		first int
	)

	/* Skip find close parentheses and result ready steps */
	if !_result.Found {
		return _result
	}

	/* Find open parentheses */
	for index := 0; index < len(*tokens); index++ {
		var _token objects.Token = (*tokens)[index]
		if _token.Type == fract.TypeOpenParenthes {
			first = index
			_result.Index = index
			_result.Found = true
			break
		}
	}

	/* Find close parentheses */
	var count int = 1
	for index := _result.Index + 1; index < len(*tokens); index++ {
		var _token objects.Token = (*tokens)[index]
		if _token.Type == fract.TypeCloseParenthes {
			count = count - 1
			if count == 0 {
				break
			}
		} else if _token.Type == fract.TypeOpenParenthes {
			count = count + 1
		}
		_result.Range = append(_result.Range, _token)
	}

	/* Remove range from original tokens */
	copy((*tokens)[first:], (*tokens)[first+len(_result.Range):])
	(*tokens)[len(*tokens)-1] = *new(objects.Token)
	*tokens = (*tokens)[:len(*tokens)-1]

	return _result
}
