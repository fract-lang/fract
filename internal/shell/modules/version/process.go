/*
	Process Function.
*/

package version

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/fract"
)

// Process Process command in module.
// command Command to process.
func Process(command string) {
	if command != "" {
		fmt.Println("This module can only be used!")
		return
	}
	fmt.Println("Fract Version [" + fract.FractVersion + "]")
}
