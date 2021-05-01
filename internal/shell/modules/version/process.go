package version

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/fract"
)

// Process command in module.
func Process(command string) {
	if command != "" {
		fmt.Println("This module can only be used!")
		return
	}
	fmt.Println("Fract Version [" + fract.FractVersion + "]")
}
