/*
	processInput Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/pkg/cli"
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// processInput Process user input.
// tokens Tokens to process.
func (i *Interpreter) processInput(tokens []obj.Token) obj.Value {
	printValue(i.processValue(&tokens))
	return obj.Value{
		Content: []string{cli.Input("")},
		Type:    fract.VALString,
	}
}
