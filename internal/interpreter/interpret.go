/*
	Interpret Function
*/

package interpreter

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
)

// Interpret Interpret file.
func (i *Interpreter) Interpret() {
	if i.Lexer.File.Path == fract.Stdin {
		// Interpret all lines.
		for i.index = 0; i.index < len(i.Tokens); i.index++ {
			i.processTokens(i.Tokens[i.index])
		}
		return
	}

	// Lexer is finished.
	if i.Lexer.Finished {
		return
	}

	{
		// Import same directory.

		path := "." + string(os.PathSeparator)
		content, err := ioutil.ReadDir(path)

		if err == nil {
			for _, current := range content {
				// Skip directories.
				if current.IsDir() || !strings.HasSuffix(current.Name(), fract.FractExtension) ||
					current.Name() == i.Lexer.File.File.Name() {
					continue
				}

				source := New(path, path+current.Name())
				source.Import()

				i.functions = append(i.functions, source.functions...)
				i.variables = append(i.variables, source.variables...)
			}
		}
	}

	i.ready()

	// Interpret all lines.
	for i.index = 0; i.index < len(i.Tokens); i.index++ {
		i.processTokens(i.Tokens[i.index])
	}
}
