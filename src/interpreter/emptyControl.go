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
	if !(*tokens).Any() {
		if i.index < i.tokenLen {
			*tokens = i.tokens.At(i.index).(*vector.Vector)
			return true
		}
		return false
	}
	i.tokens.Insert(i.index, *tokens)
	return true
}
