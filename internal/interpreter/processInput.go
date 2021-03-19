/*
	processInput Function.
*/

package interpreter

import (
	"bufio"
	"os"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// processInput Process user input.
// tokens Tokens to process.
func (i *Interpreter) processInput(tokens vector.Vector) objects.Value {
	printValue(i.processValue(&tokens))
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	value := objects.Value{
		Content: []string{},
		Type:    fract.VALString,
		Array:   true,
	}
	for _, char := range input {
		value.Content = append(value.Content, arithmetic.IntToString(char))
	}
	return value
}
