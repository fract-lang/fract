package parser

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

// Import content into destination interpeter.
func (p *Parser) Import() {
	p.ready()
	// Interpret all lines.
	for p.i = 0; p.i < len(p.Tks); p.i++ {
		switch tks := p.Tks[p.i]; tks[0].T {
		case fract.Protected: // Protected declaration.
			if len(tks) < 2 {
				fract.Error(tks[0], "Protected but what is it protected?")
			}
			second := tks[1]
			tks = tks[1:]
			if second.T == fract.Var { // Variable definition.
				p.vardec(tks, true)
			} else if second.T == fract.Func { // Function definition.
				p.funcdec(tks, true)
			} else {
				fract.Error(second, "Syntax error, you can protect only deletable objects!")
			}
		case fract.Var: // Variable definition.
			p.vardec(tks, false)
		case fract.Func: // Function definiton.
			p.funcdec(tks, false)
		case fract.Import: // Import.
			src := new(Parser)
			src.ApplyEmbedFunctions()
			src.procImport(tks)
			p.vars = append(p.vars, src.vars...)
			p.funcs = append(p.funcs, src.funcs...)
			p.Imports = append(p.Imports, src.Imports...)
		case fract.Macro: // Macro.
			p.procMacro(tks)
			if p.loopCount != -1 { // Breaked import.
				return
			}
		default:
			p.skipBlock(true)
		}
	}
}

// Information of import.
type importInfo struct {
	Name string  // Package name.
	Src  *Parser // Source of package.
}

func (p *Parser) procImport(tks []obj.Token) {
	if len(tks) == 1 {
		fract.Error(tks[0], "Imported but what?")
	}
	if tks[1].T != fract.Name && (tks[1].T != fract.Value || tks[1].Val[0] != '"' && tks[1].Val[0] != '.') {
		fract.Error(tks[1], "Import path should be string or standard path!")
	}
	j := 1
	if len(tks) > 2 {
		if tks[1].T == fract.Name {
			j = 2
		} else {
			fract.Error(tks[1], "Alias is should be name!")
		}
	}
	if j == 1 && len(tks) != 2 {
		fract.Error(tks[2], "Invalid syntax!")
	} else if j == 2 && len(tks) != 3 {
		fract.Error(tks[3], "Invalid syntax!")
	}
	src := new(Parser)
	src.ApplyEmbedFunctions()
	var imppath string
	if tks[j].T == fract.Name {
		if !strings.HasPrefix(tks[j].Val, "std") {
			fract.Error(tks[j], "Standard import should be starts with 'std' directory.")
		}
		switch tks[j].Val {
		default:
			imppath = strings.ReplaceAll(tks[j].Val, ".", string(os.PathSeparator))
		}
	} else {
		imppath = tks[0].F.P[:strings.LastIndex(tks[0].F.P, string(os.PathSeparator))+1] + p.procVal([]obj.Token{tks[j]}).D[0].String()
	}
	imppath = path.Join(fract.ExecPath, imppath)
	info, err := os.Stat(imppath)
	// Exists directory?
	if imppath != "" && (err != nil || !info.IsDir()) {
		fract.Error(tks[j], "Directory not found/access!")
	}
	infos, err := ioutil.ReadDir(imppath)
	if err != nil {
		fract.Error(tks[1], "There is a problem on import: "+err.Error())
	}
	// TODO: Improve naming.
	var name string
	if j == 1 {
		name = info.Name()
	} else {
		name = tks[1].Val
	}
	// Check name.
	for _, imp := range p.Imports {
		if imp.Name == name {
			fract.Error(tks[1], "\""+name+"\" is already defined!")
		}
	}
	for _, i := range infos {
		// Skip directories.
		if i.IsDir() || !strings.HasSuffix(i.Name(), fract.Ext) {
			continue
		}
		isrc := New(imppath + string(os.PathSeparator) + i.Name())
		isrc.loopCount = -1 //! Tag as import source.
		isrc.Import()
		src.funcs = append(src.funcs, isrc.funcs...)
		src.vars = append(src.vars, isrc.vars...)
		src.Imports = append(src.Imports, isrc.Imports...)
	}
	p.Imports = append(p.Imports,
		importInfo{
			Name: name,
			Src:  src,
		})
}
