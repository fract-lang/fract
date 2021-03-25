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
	}
	if !strings.HasSuffix(command, fract.FractExtension) {
		command += fract.FractExtension
	}
	if !fs.ExistFile(command) {
		fmt.Println("The Fract file is not exists: " + command)
		return
	}

	preter := interpreter.New(command, fract.TypeEntryFile)
	preter.ApplyEmbedFunctions()

	except.Block{
		Try: func() {
			preter.Interpret()
		},
		Catch: func(e obj.Exception) {
			os.Exit(0)
		},
	}.Do()
}
