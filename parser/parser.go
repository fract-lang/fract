package parser

import (
	"../fract"
	"../objects"
	"../utilities/list"
)

// DecomposeArithmeticProcesses Decompose and returns arithmetic processes by operators.
func DecomposeArithmeticProcesses(tokens *list.List) list.List {
	var (
		processes list.List = *list.New()
		process   objects.ArithmeticProcess
		new       bool = true
	)

	for index := 0; index < tokens.Len(); index++ {
		var _token objects.Token = tokens.Vals[index].(objects.Token)

		if _token.Type != fract.TypeOperator && _token.Type != fract.TypeValue {
			fract.Error(_token, "This is not a invalid statement!: "+_token.Value)
		}

		if new {
			new = false
			process.First = _token
			continue
		}
		if process.Operator.Value == "" {
			if _token.Type != fract.TypeOperator {
				fract.Error(_token, "Operator is not found!: "+_token.Value)
			}
			process.Operator = _token
			continue
		}
		process.Second = _token
		processes.Append(process)

		/* Reset to defaults. */
		process.First.Value = ""
		process.Second.Value = ""
		process.Operator.Value = ""
		new = true
	}

	if process.First.Value != "" {
		if process.Operator.Value == "" {
			fract.Error(tokens.Last().(objects.Token), "Operator is not found!")
		} else if process.Second.Value == "" {
			fract.Error(tokens.Last().(objects.Token), "Second value is not found!")
		}
	}

	return processes
}
