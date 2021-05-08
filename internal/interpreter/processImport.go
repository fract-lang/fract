package interpreter

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
)

func (i *Interpreter) processImport(tokens []objects.Token) {
	if len(tokens) == 1 {
		fract.Error(tokens[0], "Imported but what?")
	}

	if tokens[1].Type != fract.TypeName && (tokens[1].Type != fract.TypeValue ||
		(!strings.HasPrefix(tokens[1].Value, grammar.TokenDoubleQuote) &&
			!strings.HasPrefix(tokens[1].Value, grammar.TokenQuote))) {
		fract.Error(tokens[1], "Import path should be string or standard path!")
	}

	index := 1
	if len(tokens) > 2 {
		if tokens[1].Type == fract.TypeName {
			index = 2
		} else {
			fract.Error(tokens[1], "Alias is should be name!")
		}
	}

	if index == 1 && len(tokens) != 2 {
		fract.Error(tokens[2], "Invalid syntax!")
	} else if index == 2 && len(tokens) != 3 {
		fract.Error(tokens[3], "Invalid syntax!")
	}

	var importpath string
	if tokens[index].Type == fract.TypeName {
		if !strings.HasPrefix(tokens[index].Value, "std") {
			fract.Error(tokens[index], "Standard import should be starts with 'std' directory.")
		}
		importpath = strings.ReplaceAll(tokens[index].Value, grammar.TokenDot, string(os.PathSeparator))
	} else {
		importpath = tokens[0].File.Path[:strings.LastIndex(tokens[0].File.Path, string(os.PathSeparator))+1] +
			i.processValue([]objects.Token{tokens[index]}).Content[0].Data
	}
	
	importpath = path.Join(fract.ExecutablePath, importpath)

	info, err := os.Stat(importpath)

	// Exists directory?
	if err != nil || !info.IsDir() {
		fract.Error(tokens[index], "Directory not found/access!")
	}

	content, err := ioutil.ReadDir(importpath)
	if err != nil {
		fract.Error(tokens[1], "There is a problem on import: "+err.Error())
	}

	var name string
	if index == 1 {
		name = info.Name()
	} else {
		name = tokens[1].Value
	}

	// Check name.
	for _, _import := range i.Imports {
		if _import.Name == name {
			fract.Error(tokens[1], "'"+name+"' is already defined!")
		}
	}

	source := new(Interpreter)
	source.ApplyEmbedFunctions()
	for _, current := range content {
		// Skip directories.
		if current.IsDir() || !strings.HasSuffix(current.Name(), fract.FractExtension) {
			continue
		}

		isource := New(importpath, importpath+string(os.PathSeparator)+current.Name())
		isource.Import()

		source.functions = append(source.functions, isource.functions...)
		source.variables = append(source.variables, isource.variables...)
		source.macroDefines = append(source.macroDefines, isource.macroDefines...)
		source.Imports = append(source.Imports, isource.Imports...)
	}

	i.Imports = append(i.Imports,
		ImportInfo{
			Name:   name,
			Source: source,
		})
}
