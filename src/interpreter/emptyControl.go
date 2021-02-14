/*
	emptyControl Function.
*/

package interpreter

import (
	"../utilities/vector"
)

// emptyControl Control empty tokens and return success state.
// tokens Tokens to check.
func (i *Interpreter) emptyControl(tokens **vector.Vector) bool {
	if len((*tokens).Vals) == 0 {
		if i.index < len(i.tokens.Vals) {
			*tokens = i.tokens.Vals[i.index].(*vector.Vector)
			return true
		}
		return false
	}
	i.tokens.Insert(i.index, *tokens)
	return true
}
