/*
	VALUE PROCESSORS
*/

package interpreter

import (
	"strings"

	"../fract"
	"../fract/arithmetic"
	"../grammar"
	"../objects"
	"../parser"
	"../utilities/vector"
)

// processValue Process value.
// tokens Tokens.
func (i *Interpreter) processValue(tokens *vector.Vector) objects.Value {
	/* Check parentheses range. */
	for true {
		_range, found := parser.DecomposeBrace(tokens)

		/* Parentheses are not found! */
		if found == -1 {
			break
		}

		var _token objects.Token
		_token.Value = i.processValue(&_range).Content
		_token.Type = fract.TypeValue
		tokens.Insert(found, _token)
	}

	var value objects.Value
	value.Content = ""
	value.Type = fract.VTInteger

	// Decompose arithmetic operations.
	operations := parser.DecomposeArithmeticProcesses(tokens)

	// Process arithmetic operation.
	priorityIndex := parser.IndexProcessPriority(&operations)
	for priorityIndex != -1 {
		var operation objects.ArithmeticProcess
		operation.First = operations.At(priorityIndex - 1).(objects.Token)
		operation.Operator = operations.At(priorityIndex).(objects.Token)
		operation.Second = operations.At(priorityIndex + 1).(objects.Token)

		_token := operations.At(priorityIndex - 1).(objects.Token)
		operations.RemoveRange(priorityIndex-1, 3)
		_type, result := arithmetic.SolveArithmeticProcess(operation)
		value.Type = _type
		_token.Value = arithmetic.TypeToString(_type, result)
		operations.Insert(priorityIndex-1, _token)

		// Find next operator.
		priorityIndex = parser.IndexProcessPriority(&operations)
	}

	// Set value.
	_value, _ := arithmetic.ToFloat64(operations.First().(objects.Token).Value)
	value.Content = arithmetic.TypeToString(value.Type, _value)

	/* Set type to float if... */
	if value.Type != fract.VTFloat &&
		(strings.Index(value.Content, grammar.TokenDot) != -1 ||
			strings.Index(value.Content, grammar.TokenDot) != -1) {
		value.Type = fract.VTFloat
	}

	return value
}
