/*
	processInput Function.
*/

package interpreter

import (
	"bufio"
	"os"

	"../fract/arithmetic"
	"../objects"
	"../utilities/cli"
	"../utilities/vector"
)

// processInput Process user input.
// tokens Tokens to process.
func (i *Interpreter) processInput(tokens *vector.Vector) objects.Value {
	if len(tokens.Vals) == 0 {
		return objects.Value{
			Content: []string{cli.Input("")},
			Charray: true,
			Array:   true,
		}
	}
	printValue(i.processValue(tokens))
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
