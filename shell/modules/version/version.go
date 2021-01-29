package version

import "fmt"

// Process Process command in module.
func Process(command string) {
	if command != "" {
		fmt.Println("This module can only be used!")
		return
	}
	fmt.Println("Fract Version [0.0.1]")
}
