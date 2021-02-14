/*
	subtractBlock Function.
*/

package interpreter

import (
	"../fract"
	"../objects"
)

// subtractBlock Subtract block count.
func (i *Interpreter) subtractBlock(token *objects.Token) {
	i.blockCount--
	if i.blockCount < 0 {
		fract.Error(*token, "The extra block end defined!")
	}
}
