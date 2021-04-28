/*
	Process Function.
*/

package make

import (
	"fmt"
	"os"
	"strings"

	"github.com/fract-lang/fract/internal/interpreter"
	"github.com/fract-lang/fract/pkg/except"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/fs"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Process Process command in module.
// command Command to process.
func Process(command string) {
	if command == "" {
		fmt.Println("This module cannot only be used!")
		return
	} else if !strings.HasSuffix(command, fract.FractExtension) {
		command += fract.FractExtension
	}

	if !fs.ExistFile(command) {
		fmt.Println("The Fract file is not exists: " + command)
		return
	}

	if info, err := os.Stat("std"); err != nil || !info.IsDir() {
		fmt.Println("Standard library not found!")
		return
	}

	preter := interpreter.New(".", command)
	preter.ApplyEmbedFunctions()

	(&except.Block{
		Try: func() {
			preter.Interpret()
		},
		Catch: func(e *obj.Exception) {
			os.Exit(0)
		},
	}).Do()
}
