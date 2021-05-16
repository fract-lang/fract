package interpreter

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
)

func (i *Interpreter) Interpret() {
	if i.Lexer.File.Path == "<stdin>" {
		// Interpret all lines.
		for i.index = 0; i.index < len(i.Tokens); i.index++ {
			i.processTokens(i.Tokens[i.index])
			runtime.GC()
		}
		return
	}

	// Lexer is finished.
	if i.Lexer.Finished { return }

	i.ready()

	{
		//* Import local directory.

		dir, _ := os.Getwd()
		if pdir := path.Dir(i.Lexer.File.Path); pdir != "." {
			dir = path.Join(dir, pdir)
		}
		content, err := ioutil.ReadDir(dir)

		if err == nil {
			_, mainName := filepath.Split(i.Lexer.File.Path)
			for _, current := range content {
				// Skip directories.
				if current.IsDir() || !strings.HasSuffix(current.Name(), fract.FractExtension) || current.Name() == mainName {
					continue
				}
				
				source := New(dir, path.Join(dir, current.Name()))
				source.ApplyEmbedFunctions()
				source.Import()

				i.functions = append(i.functions, source.functions...)
				i.variables = append(i.variables, source.variables...)
				i.macroDefines = append(i.macroDefines, source.macroDefines...)
			}
		}
	}

	// Interpret all lines.
	for i.index = 0; i.index < len(i.Tokens); i.index++ {
		i.processTokens(i.Tokens[i.index])
		runtime.GC()
	}
}
