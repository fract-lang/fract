package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// DecomposeArithmeticProcesses returns arithmetic processes by operators.
func DecomposeArithmeticProcesses(tokens []objects.Token) *[]objects.Token {
	var (
		operator  bool
		lastIndex int
		processes []objects.Token
	)

	for index := 0; index < len(tokens); index++ {
		switch token := tokens[index]; token.Type {
		case fract.TypeOperator:
			if !operator {
				fract.Error(token, "Operator spam!")
			}
			lastIndex = index
			processes = append(processes, token)
			operator = false
		case fract.TypeValue, fract.TypeName, fract.TypeComma, fract.TypeBrace:
			lastIndex = index
			processes = append(processes, token)
			operator = index < len(tokens)-1
		default:
			fract.Error(token, "Invalid value!")
		}
	}

	if processes[lastIndex].Type == fract.TypeOperator {
		fract.Error(processes[len(processes)-1], "Operator defined, but for what?")
	}

	return &processes
}
