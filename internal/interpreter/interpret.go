package interpreter

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
)

func (i *Interpreter) Interpret() {
	if i.Lexer.File.Path == fract.Stdin {
		// Interpret all lines.
		for i.index = 0; i.index < len(i.Tokens); i.index++ {
			i.processTokens(i.Tokens[i.index])
			runtime.GC()
		}
		return
	}

	// Lexer is finished.
	if i.Lexer.Finished {
		return
	}

	{
		//* Import local directory.

		path := "." + string(os.PathSeparator)
		content, err := ioutil.ReadDir(path)

		if err == nil {
			mainName := i.Lexer.File.Path[strings.LastIndex(i.Lexer.File.Path, string(os.PathSeparator))+1:]
			for _, current := range content {
				// Skip directories.
				if current.IsDir() || !strings.HasSuffix(current.Name(), fract.FractExtension) ||
					current.Name() == mainName {
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
		runtime.GC()
	}
}
