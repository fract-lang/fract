/*
	Input Function.
*/

package cli

import (
	"fmt"
)

// Input Returns input from command-line.
// message Input message.
func Input(message string) string {
	fmt.Print(message)
	var input string
	fmt.Sprintln(&input)
	return input
}
