/*
	Process Function.
*/

package make

import (
	"fmt"
	"strings"

	"../../../fract"
	"../../../interpreter"
	"../../../utilities/fs"
)

// Process Process command in module.
// command Command to process.
func Process(command string) {
	if command == "" {
		fmt.Println("This module cannot only be used!")
		return
	}
	if !strings.HasSuffix(command, fract.FractExtension) {
		command += fract.FractExtension
	}
	if !fs.ExistFile(command) {
		fmt.Println("The Fract file is not exists: " + command)
		return
	}

	preter := interpreter.New(command, fract.TypeEntryFile)
	preter.Interpret()
}
