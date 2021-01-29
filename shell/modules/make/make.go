package make

import (
	"fmt"
	"strings"

	"../../../utilities/fs"
)

// Process Process command in module.
func Process(command string) {
	if command == "" {
		fmt.Println("This module cannot only be used!")
		return
	}
	if strings.HasSuffix(command, ".fract") {
		command += ".fract"
	}
	if !fs.ExistsFile(command) {
		fmt.Println("The Fract file is not exists: " + command)
		return
	}
	/* Parser commands */
	fmt.Println("Success!")
}
