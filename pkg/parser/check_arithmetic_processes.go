package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// CheckArithmeticProcesses returns arithmetic processes by operators.
func CheckArithmeticProcesses(tokens []objects.Token) {
	var (
		operator  bool
		lastIndex int
	)
	for index := 0; index < len(tokens); index++ {
		switch token := tokens[index]; token.Type {
		case fract.TypeOperator:
			if !operator {
				fract.Error(token, "Operator spam!")
			}
			lastIndex = index
			operator = false
		case fract.TypeValue, fract.TypeName, fract.TypeComma, fract.TypeBrace:
			lastIndex = index
			operator = index < len(tokens)-1
		default:
			fract.Error(token, "Invalid value!")
		}
	}
	if lastIndex < len(tokens)-1 {
		token := tokens[lastIndex]
		if token.Type == fract.TypeOperator && !operator {
			fract.Error(token, "Operator spam!")
		} else if token.Type != fract.TypeValue && token.Type != fract.TypeName &&
			token.Type != fract.TypeBrace && token.Type != fract.TypeComma {
			fract.Error(token, "Invalid value!")
		}
	}
	if tokens[len(tokens)-1].Type == fract.TypeOperator {
		fract.Error(tokens[len(tokens)-1], "Operator overflow!")
	}
}
