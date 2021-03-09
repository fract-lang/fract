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
func IndexBlockDeclare(tokens vector.Vector) int {
	for index := range tokens.Vals {
		if tokens.Vals[index].(objects.Token).Type == fract.TypeBlock {
			return index
		}
	}
	return -1
}
