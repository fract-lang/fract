package interpreter

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// Import content into destination interpeter.
func (i *Interpreter) Import() {
	i.ready()
	// Interpret all lines.
	for i.index = 0; i.index < len(i.Tokens); i.index++ {
		switch tokens := i.Tokens[i.index]; tokens[0].Type {
		case fract.TypeProtected: // Protected declaration.
			if len(tokens) < 2 {
				fract.Error(tokens[0], "Protected but what is it protected?")
			}
			second := tokens[1]
			tokens = tokens[1:]
			if second.Type == fract.TypeVariable { // Variable definition.
				i.processVariableDeclaration(tokens, true)
			} else if second.Type == fract.TypeFunction { // Function definition.
				i.processFunctionDeclaration(tokens, true)
			} else {
				fract.Error(second, "Syntax error, you can protect only deletable objects!")
			}
		case fract.TypeVariable: // Variable definition.
			i.processVariableDeclaration(tokens, false)
		case fract.TypeFunction: // Function definiton.
			i.processFunctionDeclaration(tokens, false)
		case fract.TypeImport: // Import.
			source := new(Interpreter)
			source.ApplyEmbedFunctions()
			source.processImport(tokens)

			i.variables = append(i.variables, source.variables...)
			i.functions = append(i.functions, source.functions[:]...)
			i.Imports = append(i.Imports, source.Imports...)
		case fract.TypeMacro: // Macro.
			i.processMacro(tokens)
			if i.loopCount != -1 { // Breaked import.
				return
			}
		default:
			i.skipBlock(true)
		}
	}
}

// Information of import.
type importInfo struct {
	Name   string       // Package name.
	Source *Interpreter // Source of package.
}

func (i *Interpreter) processImport(tokens []objects.Token) {
	if len(tokens) == 1 {
		fract.Error(tokens[0], "Imported but what?")
	}
	if tokens[1].Type != fract.TypeName && (tokens[1].Type != fract.TypeValue || tokens[1].Value[0] != '"' && tokens[1].Value[0] != '.') {
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
	source := new(Interpreter)
	source.ApplyEmbedFunctions()
	var importpath string
	if tokens[index].Type == fract.TypeName {
		if !strings.HasPrefix(tokens[index].Value, "std") {
			fract.Error(tokens[index], "Standard import should be starts with 'std' directory.")
		}
		switch tokens[index].Value {
		default:
			importpath = strings.ReplaceAll(tokens[index].Value, ".", string(os.PathSeparator))
		}
	} else {
		importpath = tokens[0].File.Path[:strings.LastIndex(tokens[0].File.Path, string(os.PathSeparator))+1] +
			i.processValue([]objects.Token{tokens[index]}).Content[0].String()
	}
	importpath = path.Join(fract.ExecutablePath, importpath)
	info, err := os.Stat(importpath)
	// Exists directory?
	if importpath != "" && (err != nil || !info.IsDir()) {
		fract.Error(tokens[index], "Directory not found/access!")
	}
	content, err := ioutil.ReadDir(importpath)
	if err != nil {
		fract.Error(tokens[1], "There is a problem on import: "+err.Error())
	}
	// TODO: Improve naming.
	var name string
	if index == 1 {
		name = info.Name()
	} else {
		name = tokens[1].Value
	}
	// Check name.
	for _, _import := range i.Imports {
		if _import.Name == name {
			fract.Error(tokens[1], "\""+name+"\" is already defined!")
		}
	}
	for _, current := range content {
		// Skip directories.
		if current.IsDir() || !strings.HasSuffix(current.Name(), fract.FractExtension) {
			continue
		}

		isource := New(importpath, importpath+string(os.PathSeparator)+current.Name())
		isource.loopCount = -1 //! Tag as import source.
		isource.Import()

		source.functions = append(source.functions, isource.functions...)
		source.variables = append(source.variables, isource.variables...)
		source.macroDefines = append(source.macroDefines, isource.macroDefines...)
		source.Imports = append(source.Imports, isource.Imports...)
	}
	i.Imports = append(i.Imports,
		importInfo{
			Name:   name,
			Source: source,
		})
}
