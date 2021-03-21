/*
	processInput Function.
*/

package interpreter

import (
	"bufio"
	"os"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// processInput Process user input.
// tokens Tokens to process.
func (i *Interpreter) processInput(tokens []obj.Token) obj.Value {
	printValue(i.processValue(&tokens))
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	value := obj.Value{
		Content: []string{},
		Type:    fract.VALString,
		Array:   true,
	}
	for _, char := range input {
		value.Content = append(value.Content, arithmetic.IntToString(char))
	}
	return value
}
