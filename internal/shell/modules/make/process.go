/*
	Process Function.
*/

package make

import (
	"fmt"
	"strings"

	"github.com/fract-lang/fract/internal/interpreter"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/fs"
)

// Process Process command in module.
// command Command to process.
func Process(command string) {
	if command == "" {
		fmt.Println("This module can not only be used!")
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
