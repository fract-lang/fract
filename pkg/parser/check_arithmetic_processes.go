package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// Check arithmetic processes validity.
func CheckArithmeticProcesses(tokens []objects.Token) {
	var (
		operator bool
		brace    int
	)
	for index := 0; index < len(tokens); index++ {
		switch token := tokens[index]; token.Type {
		case fract.TypeOperator:
			if !operator {
				fract.Error(token, "Operator spam!")
			}
			operator = false
		case fract.TypeValue, fract.TypeName, fract.TypeComma, fract.TypeBrace:
			switch token.Type {
			case fract.TypeBrace:
				if token.Value == "(" || token.Value == "[" || token.Value == "{" {
					brace++
				} else {
					brace--
				}
			case fract.TypeComma:
				if brace == 0 {
					fract.Error(token, "Invalid syntax!")
				}
			}
			operator = index < len(tokens)-1
		default:
			fract.Error(token, "Invalid syntax!")
		}
	}
	if tokens[len(tokens)-1].Type == fract.TypeOperator {
		fract.Error(tokens[len(tokens)-1], "Operator overflow!")
	}
}
