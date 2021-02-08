/*
	IndexBlockDeclare Function.
*/

package parser

import (
	"../fract"
	"../objects"
	"../utilities/vector"
)

// IndexBlockDeclare Find index of block declare.
// tokens Tokens to search.
func IndexBlockDeclare(tokens *vector.Vector) int {
	for index := 0; index < tokens.Len(); index++ {
		current := tokens.At(index).(objects.Token)
		if current.Type == fract.TypeBlock {
			return index
		}
	}
	return -1
}
