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

	for index, token := range tokens {
		if token.Type == fract.TypeOperator {
			if !operator {
				fract.Error(token, "Operator spam!")
			}
			lastIndex = index
			processes = append(processes, token)
			operator = false
		} else if token.Type == fract.TypeValue || token.Type == fract.TypeName ||
			token.Type == fract.TypeBooleanTrue || token.Type == fract.TypeBooleanFalse ||
			token.Type == fract.TypeBrace || token.Type == fract.TypeComma {
			lastIndex = index
			processes = append(processes, token)
			operator = index < len(tokens)-1
		} else {
			fract.Error(token, "Invalid value!")
		}
	}

	if lastIndex < len(tokens) {
		token := tokens[lastIndex]
		if token.Type == fract.TypeOperator {
			if !operator {
				fract.Error(token, "Operator spam!")
			}
			processes = append(processes, token)
		} else if token.Type == fract.TypeValue || token.Type == fract.TypeName ||
			token.Type == fract.TypeBooleanTrue || token.Type == fract.TypeBooleanFalse ||
			token.Type == fract.TypeBrace || token.Type == fract.TypeComma {
			processes = append(processes, token)
		} else {
			fract.Error(token, "Invalid value!")
		}
	}

	if processes[len(processes)-1].Type == fract.TypeOperator {
		fract.Error(processes[len(processes)-1], "Operator defined, but for what?")
	}

	return &processes
}
