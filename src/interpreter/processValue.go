/*
	processValue Function
*/

package interpreter

import (
	"strings"

	"../fract"
	"../fract/arithmetic"
	"../fract/dt"
	"../fract/name"
	"../grammar"
	"../objects"
	"../parser"
	"../utilities/vector"
)

// processVariableName Process value is variable name?
// token Token to process.
// operations All operations.
// index Index of token.
func (i *Interpreter) processVariableName(token *objects.Token,
	operations *vector.Vector, index int) int {
	if token.Type == fract.TypeName {
		vindex := name.VarIndexByName(i.vars, token.Value)
		if vindex == -1 {
			fract.Error(*token, "Name is not defined!: "+token.Value)
		}
		if index < operations.Len()-1 {
			next := operations.At(index + 1).(objects.Token)
			// Array?
			if next.Type == fract.TypeBrace && next.Value == grammar.TokenLBracket {
				// Find close bracket.
				cindex := index + 1
				for ; cindex < operations.Len(); cindex++ {
					current := operations.At(cindex).(objects.Token)
					if current.Type == fract.TypeBrace && current.Value == grammar.TokenRBracket {
						break
					}
				}

				valueList := operations.Sublist(index+2, cindex-index-2)
				position, err := arithmetic.ToInt64(i.processValue(&valueList).Content[0])
				if err != nil {
					fract.Error(operations.At(cindex).(objects.Token), "Value out of range!")
				}
				variable := i.vars.At(vindex).(objects.Variable)
				if position < 0 || position >= int64(len(variable.Value)) {
					fract.Error(operations.At(cindex).(objects.Token), "Index is out of range!")
				}
				operations.RemoveRange(index, cindex-index)
				token.Value = variable.Value[position]
				return valueList.Len() - 1
			}
		}

		token.Value = i.vars.At(vindex).(objects.Variable).Value[0]
	} else if token.Type == fract.TypeBrace && token.Value == grammar.TokenRBracket {
		// Find close bracket.
		oindex := index - 1
		for ; oindex >= 0; oindex-- {
			current := operations.At(oindex).(objects.Token)
			if current.Type == fract.TypeBrace && current.Value == grammar.TokenLBracket {
				break
			}
		}
		// Finished?
		if oindex == 0 {
			fract.Error(operations.First().(objects.Token), "Index error!")
		}

		nameToken := operations.At(oindex - 1).(objects.Token)
		vindex := name.VarIndexByName(i.vars, nameToken.Value)
		if vindex == -1 {
			fract.Error(*token, "Name is not defined!: "+nameToken.Value)
		}
		valueList := operations.Sublist(oindex+1, index-oindex-1)
		position, err := arithmetic.ToInt64(i.processValue(&valueList).Content[0])
		if err != nil {
			fract.Error(operations.At(oindex).(objects.Token), "Value out of range!")
		}
		variable := i.vars.At(vindex).(objects.Variable)
		if position < 0 || position >= int64(len(variable.Value)) {
			fract.Error(operations.At(oindex).(objects.Token), "Index is out of range!")
		}
		operations.RemoveRange(oindex-1, index-oindex+1)
		token.Value = variable.Value[position]
		return index - oindex + 1
	}
	return 0
}

// processValue Process value.
// tokens Tokens.
func (i *Interpreter) processValue(tokens *vector.Vector) objects.Value {
	// Is array expression?
	first := tokens.First().(objects.Token)
	if first.Type == fract.TypeBrace && (first.Value == grammar.TokenLBrace ||
		first.Value == grammar.TokenLBracket) {
		return i.processArrayValue(tokens)
	}

	/* Check parentheses range. */
	for true {
		_range, found := parser.DecomposeBrace(tokens, grammar.TokenLParenthes,
			grammar.TokenRParenthes)

		/* Parentheses are not found! */
		if found == -1 {
			break
		}

		var _token objects.Token
		_token.Value = i.processValue(&_range).Content[0]
		_token.Type = fract.TypeValue
		tokens.Insert(found, _token)
	}

	var value objects.Value
	value.Type = fract.VTInteger

	// Is conditional expression?
	if i.isConditional(tokens) {
		value.Content = []string{arithmetic.IntToString(i.processCondition(tokens))}
		return value
	}

	// Decompose arithmetic operations.
	operations := parser.DecomposeArithmeticProcesses(tokens)

	// Process arithmetic operation.
	priorityIndex := parser.IndexProcessPriority(&operations)
	looped := priorityIndex != -1
	for priorityIndex != -1 {
		var operation objects.ArithmeticProcess
		operation.First = operations.At(priorityIndex - 1).(objects.Token)
		// First value is a name?
		priorityIndex -= i.processVariableName(&operation.First, &operations, priorityIndex-1)

		operation.Operator = operations.At(priorityIndex).(objects.Token)

		operation.Second = operations.At(priorityIndex + 1).(objects.Token)
		// Second value is a name?
		priorityIndex -= i.processVariableName(&operation.Second, &operations, priorityIndex+1)

		_token := operations.At(priorityIndex - 1).(objects.Token)
		operations.RemoveRange(priorityIndex-1, 3)
		_type, result := arithmetic.SolveArithmeticProcess(operation)
		value.Type = _type
		_token.Value = arithmetic.TypeToString(_type, result)
		_token.Type = fract.TypeValue
		operations.Insert(priorityIndex-1, _token)

		// Find next operator.
		priorityIndex = parser.IndexProcessPriority(&operations)
	}

	// Set value.
	first = operations.First().(objects.Token)

	// First value is a name?
	if first.Type == fract.TypeName && !looped {
		index := name.VarIndexByName(i.vars, first.Value)
		if index == -1 {
			fract.Error(first,
				"Name is not defined!: "+first.Value)
		}
		variable := i.vars.At(index).(objects.Variable)
		// Is Array?
		if variable.Array && operations.Len() == 1 {
			value.Content = variable.Value
			if dt.IsFloatType(variable.Type) {
				value.Type = fract.VTFloatArray
			} else {
				value.Type = fract.VTIntegerArray
			}
			return value
		}
		i.processVariableName(&first, &operations, 0)
		value.Content = []string{first.Value}
		return value
	}

	_value, err := arithmetic.ToFloat64(first.Value)
	if err != nil {
		fract.Error(first, "Value out of range!")
	}
	if arithmetic.IsFloatValue(first.Value) {
		value.Type = fract.VTFloat
	}
	value.Content = []string{arithmetic.TypeToString(value.Type, _value)}

	/* Set type to float if... */
	if value.Type != fract.VTFloat &&
		(strings.Index(value.Content[0], grammar.TokenDot) != -1 ||
			strings.Index(value.Content[0], grammar.TokenDot) != -1) {
		value.Type = fract.VTFloat
	}

	return value
}
