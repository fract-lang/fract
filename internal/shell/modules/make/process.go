package make

import (
	"fmt"
	"os"
	"strings"

	"github.com/fract-lang/fract/internal/interpreter"
	"github.com/fract-lang/fract/pkg/except"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/fs"
	"github.com/fract-lang/fract/pkg/objects"
)

// Process command in module.
func Process(command string) {
	if command == "" {
		fmt.Println("This module cannot only be used!")
		return
	} else if !strings.HasSuffix(command, fract.FractExtension) {
		command += fract.FractExtension
	}
	if !fs.ExistsFile(command) {
		fmt.Println("The Fract file is not exists: " + command)
		return
	}

	preter := interpreter.New(".", command)
	preter.ApplyEmbedFunctions()
	(&except.Block{
		Try: preter.Interpret,
		Catch: func(e *objects.Exception) {
			os.Exit(0)
		},
	}).Do()
}
