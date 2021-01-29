package exit

import (
	"fmt"
	"os"
)

// Process Process command in module.
func Process(command string) {
	if command != "" {
		fmt.Println("This module can only be used!")
		return
	}
	os.Exit(0)
}
