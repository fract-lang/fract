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
				fract.IPanic(tks[0], obj.SyntaxPanic, "Define is not given!")
			}
			second := tks[1]
			tks = tks[1:]
			switch second.T {
			case fract.Var: // Variable definition.
				p.vardec(tks, true)
			case fract.Func: // Function definition.
				p.funcdec(tks, true)
			default:
				fract.IPanic(second, obj.SyntaxPanic, "Can protect only deletable objects!")
			}
		case fract.Var: // Variable definition.
			p.vardec(tks, false)
		case fract.Func: // Function definiton.
			p.funcdec(tks, false)
		case fract.Import: // Import.
			src := new(Parser)
			src.AddBuiltInFuncs()
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

func (p *Parser) procImport(tks obj.Tokens) {
	if len(tks) == 1 {
		fract.IPanic(tks[0], obj.SyntaxPanic, "Import path is not given!")
	}
	if tks[1].T != fract.Name && (tks[1].T != fract.Value || tks[1].V[0] != '"' && tks[1].V[0] != '.') {
		fract.IPanic(tks[1], obj.ValuePanic, "Import path should be string or standard path!")
	}
	j := 1
	if len(tks) > 2 {
		if tks[1].T == fract.Name {
			j = 2
		} else {
			fract.IPanic(tks[1], obj.NamePanic, "Alias is should be a invalid name!")
		}
	}
	if j == 1 && len(tks) != 2 {
		fract.IPanic(tks[2], obj.SyntaxPanic, "Invalid syntax!")
	} else if j == 2 && len(tks) != 3 {
		fract.IPanic(tks[3], obj.SyntaxPanic, "Invalid syntax!")
	}
	src := new(Parser)
	src.AddBuiltInFuncs()
	var imppath string
	if tks[j].T == fract.Name {
		if !strings.HasPrefix(tks[j].V, "std") {
			fract.IPanic(tks[j], obj.ValuePanic, "Standard import should be starts with 'std' directory.")
		}
		switch tks[j].V {
		default:
			imppath = strings.ReplaceAll(tks[j].V, ".", string(os.PathSeparator))
		}
	} else {
		imppath = tks[0].F.P[:strings.LastIndex(tks[0].F.P, string(os.PathSeparator))+1] + p.procVal(obj.Tokens{tks[j]}).D[0].String()
	}
	imppath = path.Join(fract.ExecPath, imppath)
	info, err := os.Stat(imppath)
	// Exists directory?
	if imppath != "" && (err != nil || !info.IsDir()) {
		fract.IPanic(tks[j], obj.PlainPanic, "Directory not found/access!")
	}
	infos, err := ioutil.ReadDir(imppath)
	if err != nil {
		fract.IPanic(tks[1], obj.PlainPanic, "There is a problem on import: "+err.Error())
	}
	// TODO: Improve naming.
	var name string
	if j == 1 {
		name = info.Name()
	} else {
		name = tks[1].V
	}
	// Check name.
	for _, imp := range p.Imports {
		if imp.Name == name {
			fract.IPanic(tks[1], obj.NamePanic, "\""+name+"\" is already defined!")
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
