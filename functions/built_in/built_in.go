package built_in

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
	"github.com/fract-lang/fract/pkg/value"
)

// Exit from application with code.
func Exit(tk obj.Token, args []obj.Var) {
	c := args[0].V
	if c.T != value.Single || c.D[0].T != value.Int {
		fract.Panic(tk, obj.ValuePanic, "Exit code is only be integer!")
	}
	ec, _ := strconv.ParseInt(c.D[0].String(), 10, 64)
	os.Exit(int(ec))
}

// Float convert object to float.
func Float(parameters []obj.Var) value.Val {
	return value.Val{D: []value.Data{{
		D: fmt.Sprintf(fract.FloatFormat, value.Conv(parameters[0].V.D[0].String())),
		T: value.Float,
	}}}
}

// Input returns input from command-line.
func Input(args []obj.Var) value.Val {
	args[0].V.Print()
	//! Don't use fmt.Scanln
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	return value.Val{
		D: []value.Data{{D: s.Text(), T: value.Str}},
	}
}

// Int convert object to integer.
func Int(args []obj.Var) value.Val {
	switch args[1].V.D[0].D { // Cast type.
	case "strcode":
		var v value.Val
		for _, byt := range []byte(args[0].V.D[0].String()) {
			v.D = append(v.D, value.Data{D: fmt.Sprint(byt), T: value.Int})
		}
		v.T = value.Array
		return v
	default: // Object.
		return value.Val{
			D: []value.Data{{
				D: fmt.Sprint(int(value.Conv(args[0].V.D[0].String()))),
				T: value.Int,
			}},
		}
	}
}

// Len returns length of object.
func Len(args []obj.Var) value.Val {
	arg := args[0].V
	if arg.T == value.Array {
		return value.Val{D: []value.Data{{D: fmt.Sprint(len(arg.D))}}}
	} else if arg.T == value.Map {
		return value.Val{D: []value.Data{{D: fmt.Sprint(len(arg.D[0].D.(map[interface{}]value.Val)))}}}
	} else if arg.D[0].T == value.Str {
		return value.Val{D: []value.Data{{D: fmt.Sprint(len(arg.D[0].String()))}}}
	}
	return value.Val{D: []value.Data{{D: "0"}}}
}

// Calloc array by size.
func Calloc(tk obj.Token, args []obj.Var) value.Val {
	sz := args[0].V
	if sz.T != value.Single || sz.D[0].T != value.Int {
		fract.Panic(tk, obj.ValuePanic, "Size is only be integer!")
	}
	szv, _ := strconv.Atoi(sz.D[0].String())
	if szv < 0 {
		fract.Panic(tk, obj.ValuePanic, "Size should be minimum zero!")
	}
	v := value.Val{T: value.Array}
	if szv > 0 {
		var index int
		for ; index < szv; index++ {
			v.D = append(v.D, value.Data{D: "0", T: value.Int})
		}
	} else {
		v.D = []value.Data{}
	}
	return v
}

// Realloc array by size.
func Realloc(tk obj.Token, args []obj.Var) value.Val {
	szv, _ := strconv.Atoi(args[1].V.D[0].String())
	if szv < 0 {
		fract.Panic(tk, obj.ValuePanic, "Size should be minimum zero!")
	}
	var (
		b = args[0].V.D
		v = value.Val{T: value.Array}
		c = 0
	)
	if len(b) <= szv {
		v.D = b
		c = len(b)
	} else {
		v.D = b[:szv]
		return v
	}
	for ; c <= szv; c++ {
		v.D = append(v.D, value.Data{D: "0", T: value.Int})
	}
	return v
}

// Print values to cli.
func Print(tk obj.Token, args []obj.Var) {
	if args[0].V.D == nil {
		fract.Panic(tk, obj.ValuePanic, "Value is not printable!")
	}
	for _, d := range args[0].V.D {
		fmt.Print(d)
	}
}

// Print values to cli with new line.
func Println(tk obj.Token, args []obj.Var) {
	Print(tk, args)
	println()
}

// Range returns array by parameters.
func Range(tk obj.Token, args []obj.Var) value.Val {
	start := args[0].V
	to := args[1].V
	step := args[2].V
	if start.T != value.Single {
		fract.Panic(tk, obj.ValuePanic, "\"start\" argument should be numeric!")
	} else if to.T != value.Single {
		fract.Panic(tk, obj.ValuePanic, "\"to\" argument should be numeric!")
	} else if step.T != value.Single {
		fract.Panic(tk, obj.ValuePanic, "\"step\" argument should be numeric!")
	}
	if start.D[0].T != value.Int && start.D[0].T != value.Float || to.D[0].T != value.Int &&
		to.D[0].T != value.Float || step.D[0].T != value.Int && step.D[0].T != value.Float {
		fract.Panic(tk, obj.ValuePanic, "Values should be integer or float!")
	}
	startV, _ := strconv.ParseFloat(start.D[0].String(), 64)
	toV, _ := strconv.ParseFloat(to.D[0].String(), 64)
	stepV, _ := strconv.ParseFloat(step.D[0].String(), 64)
	if stepV <= 0 {
		return value.Val{D: nil, T: value.Array}
	}
	var t uint8
	if start.D[0].T == value.Float || to.D[0].T == value.Float || step.D[0].T == value.Float {
		t = value.Float
	}
	rv := value.Val{D: []value.Data{}, T: value.Array}
	if startV <= toV {
		for ; startV <= toV; startV += stepV {
			d := value.Data{D: fmt.Sprintf(fract.FloatFormat, startV), T: t}
			d.D = d.Format()
			rv.D = append(rv.D, d)
		}
	} else {
		for ; startV >= toV; startV -= stepV {
			d := value.Data{D: fmt.Sprintf(fract.FloatFormat, startV), T: t}
			d.D = d.Format()
			rv.D = append(rv.D, d)
		}
	}
	return rv
}

// String convert object to string.
func String(args []obj.Var) value.Val {
	switch args[1].V.D[0].D {
	case "parse":
		str := ""
		if val := args[0].V; val.T == value.Array {
			if len(val.D) == 0 {
				str = "[]"
			} else {
				var sb strings.Builder
				sb.WriteByte('[')
				for _, data := range val.D {
					sb.WriteString(data.String() + " ")
				}
				str = sb.String()[:sb.Len()-1] + "]"
			}
		} else {
			str = args[0].V.D[0].String()
		}
		return value.Val{D: []value.Data{{D: str, T: value.Str}}}
	case "bytecode":
		v := args[0].V
		var sb strings.Builder
		for _, d := range v.D {
			if d.T != value.Int {
				sb.WriteByte(' ')
			}
			r, _ := strconv.ParseInt(d.String(), 10, 32)
			sb.WriteByte(byte(r))
		}
		return value.Val{D: []value.Data{{D: sb.String(), T: value.Str}}}
	default: // Object.
		return value.Val{D: []value.Data{{D: fmt.Sprint(args[0].V), T: value.Str}}}
	}
}

// Append source values to destination array.
func Append(tk obj.Token, args []obj.Var) value.Val {
	src := args[0].V
	if src.T != value.Array {
		fract.Panic(tk, obj.ValuePanic, "\"src\" must be array!")
	}
	for _, d := range args[1].V.D {
		src.D = append(src.D, value.Data{D: d.D, T: d.T})
	}
	return src
}
