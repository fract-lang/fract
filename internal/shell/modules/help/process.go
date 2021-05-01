package help

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/cli"
)

// Process command in module.
func Process(command string) {
	if command != "" {
		fmt.Println("This module can only be used!")
		return
	}
	cli.PrintMapAsTable(map[string]string{
		"make":    "Interprete Fract code.",
		"version": "Show version.",
		"help":    "Show help.",
		"exit":    "Exit.",
	})
}
