/*
	processImport Function.
*/

package interpreter

import (
	"os"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
)

// processImport Process import.
// tokens Tokens to process.
func (i *Interpreter) processImport(tokens []objects.Token) {
	if len(tokens) == 1 {
		fract.Error(tokens[0], "Imported but what?")
	}

	if tokens[1].Type != fract.TypeValue ||
		(!strings.HasPrefix(tokens[1].Value, grammar.TokenDoubleQuote) &&
			!strings.HasPrefix(tokens[1].Value, grammar.TokenQuote)) {
		fract.Error(tokens[1], "Import path should be string!")
	}

	valueList := tokens[1:]
	path := tokens[0].File.Path[:strings.LastIndex(tokens[0].File.Path, string(os.PathSeparator))+1] +
		i.processValue(&valueList).Content[0].Data

	info, err := os.Stat(path)

	// Exists directory?
	if err != nil || !info.IsDir() {
		fract.Error(tokens[1], "Directory not found/access!")
	}

	content, err := os.ReadDir(path)
	if err != nil {
		fract.Error(tokens[1], "There is a problem on import: "+err.Error())
	}

	for _, current := range content {
		// Skip directories.
		if current.IsDir() {
			continue
		}

		New(path, path+string(os.PathSeparator)+current.Name()).Import(i, info.Name())
	}
}
