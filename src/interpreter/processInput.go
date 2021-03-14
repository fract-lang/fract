/*
	processInput Function.
*/

package interpreter

import (
	"bufio"
	"os"

	"github.com/fract-lang/src/fract/arithmetic"
	"github.com/fract-lang/src/objects"
	"github.com/fract-lang/src/utils/vector"
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
		Charray: true,
		Array:   true,
	}
	for index := range input {
		value.Content = append(value.Content, arithmetic.IntToString(input[index]))
	}
	return value
}
