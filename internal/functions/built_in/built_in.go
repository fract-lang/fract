package built_in

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

// Exit from application with code.
func Exit(tk obj.Token, args []obj.Var) {
	c := args[0].V
	if c.Arr {
		fract.Panic(tk, obj.ValuePanic, "Array is not a valid value!")
	} else if c.D[0].T != obj.VInt {
		fract.Panic(tk, obj.ValuePanic, "Exit code is only be integer!")
	}
	ec, _ := strconv.ParseInt(c.D[0].String(), 10, 64)
	os.Exit(int(ec))
}

// Float convert object to float.
func Float(parameters []obj.Var) obj.Value {
	return obj.Value{D: []obj.Data{
		{D: fmt.Sprintf(fract.FloatFormat, arithmetic.Value(parameters[0].V.D[0].String())), T: obj.VFloat},
	}}
}

// Input returns input from command-line.
func Input(args []obj.Var) obj.Value {
	args[0].V.Print()
	//! Don't use fmt.Scanln
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	return obj.Value{
		D: []obj.Data{{D: s.Text(), T: obj.VStr}},
	}
}

// Int convert object to integer.
func Int(args []obj.Var) obj.Value {
	switch args[1].V.D[0].D { // Cast type.
	case "strcode":
		var v obj.Value
		for _, byt := range []byte(args[0].V.D[0].String()) {
			v.D = append(v.D, obj.Data{D: fmt.Sprint(byt), T: obj.VInt})
		}
		v.Arr = len(v.D) > 1
		return v
	default: // Object.
		return obj.Value{
			D: []obj.Data{{D: fmt.Sprint(int(arithmetic.Value(args[0].V.D[0].String()))), T: obj.VInt}},
		}
	}
}

// Len returns length of object.
func Len(args []obj.Var) obj.Value {
	arg := args[0].V
	if arg.Arr {
		return obj.Value{D: []obj.Data{{D: fmt.Sprint(len(arg.D))}}}
	} else if arg.D[0].T == obj.VStr {
		return obj.Value{D: []obj.Data{{D: fmt.Sprint(len(arg.D[0].String()))}}}
	}
	return obj.Value{D: []obj.Data{{D: "0"}}}
}

// Calloc array by size.
func Calloc(tk obj.Token, args []obj.Var) obj.Value {
	sz := args[0].V
	if sz.Arr {
		fract.Panic(tk, obj.ValuePanic, "Array is not a valid value!")
	} else if sz.D[0].T != obj.VInt {
		fract.Panic(tk, obj.ValuePanic, "Size is only be integer!")
	}
	szv, _ := strconv.Atoi(sz.D[0].String())
	if szv < 0 {
		fract.Panic(tk, obj.ValuePanic, "Size should be minimum zero!")
	}
	v := obj.Value{Arr: true}
	if szv > 0 {
		var index int
		for ; index < szv; index++ {
			v.D = append(v.D, obj.Data{D: "0", T: obj.VInt})
		}
	} else {
		v.D = []obj.Data{}
	}
	return v
}

// Realloc array by size.
func Realloc(tk obj.Token, args []obj.Var) obj.Value {
	szv, _ := strconv.Atoi(args[1].V.D[0].String())
	if szv < 0 {
		fract.Panic(tk, obj.ValuePanic, "Size should be minimum zero!")
	}
	var (
		b = args[0].V.D
		v = obj.Value{Arr: true}
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
		v.D = append(v.D, obj.Data{D: "0", T: obj.VInt})
	}
	return v
}

// Print values to cli.
func Print(tk obj.Token, args []obj.Var) {
	if args[0].V.D == nil {
		fract.Panic(tk, obj.ValuePanic, "Value is not printable!")
	}
	args[0].V.Print()
	args[1].V.Print()
}

// Range returns array by parameters.
func Range(tk obj.Token, args []obj.Var) obj.Value {
	start := args[0].V
	to := args[1].V
	step := args[2].V
	if start.Arr {
		fract.Panic(tk, obj.ValuePanic, "\"start\" argument should be numeric!")
	} else if to.Arr {
		fract.Panic(tk, obj.ValuePanic, "\"to\" argument should be numeric!")
	} else if step.Arr {
		fract.Panic(tk, obj.ValuePanic, "\"step\" argument should be numeric!")
	}
	if start.D[0].T != obj.VInt && start.D[0].T != obj.VFloat || to.D[0].T != obj.VInt &&
		to.D[0].T != obj.VFloat || step.D[0].T != obj.VInt && step.D[0].T != obj.VFloat {
		fract.Panic(tk, obj.ValuePanic, "Values should be integer or float!")
	}
	startV, _ := strconv.ParseFloat(start.D[0].String(), 64)
	toV, _ := strconv.ParseFloat(to.D[0].String(), 64)
	stepV, _ := strconv.ParseFloat(step.D[0].String(), 64)
	if stepV <= 0 {
		return obj.Value{
			D:   nil,
			Arr: true,
		}
	}
	var t uint8
	if start.D[0].T == obj.VFloat || to.D[0].T == obj.VFloat || step.D[0].T == obj.VFloat {
		t = obj.VFloat
	}
	rv := obj.Value{
		D:   []obj.Data{},
		Arr: true,
	}
	if startV <= toV {
		for ; startV <= toV; startV += stepV {
			d := obj.Data{D: fmt.Sprintf(fract.FloatFormat, startV), T: t}
			d.D = d.Format()
			rv.D = append(rv.D, d)
		}
	} else {
		for ; startV >= toV; startV -= stepV {
			d := obj.Data{D: fmt.Sprintf(fract.FloatFormat, startV), T: t}
			d.D = d.Format()
			rv.D = append(rv.D, d)
		}
	}
	return rv
}

// String convert object to string.
func String(args []obj.Var) obj.Value {
	switch args[1].V.D[0].D {
	case "parse":
		str := ""
		if value := args[0].V; value.Arr {
			if len(value.D) == 0 {
				str = "[]"
			} else {
				var sb strings.Builder
				sb.WriteByte('[')
				for _, data := range value.D {
					sb.WriteString(data.String() + " ")
				}
				str = sb.String()[:sb.Len()-1] + "]"
			}
		} else {
			str = args[0].V.D[0].String()
		}
		return obj.Value{
			D: []obj.Data{{D: str, T: obj.VStr}},
		}
	case "bytecode":
		v := args[0].V
		var sb strings.Builder
		for _, d := range v.D {
			if d.T != obj.VInt {
				sb.WriteByte(' ')
			}
			r, _ := strconv.ParseInt(d.String(), 10, 32)
			sb.WriteByte(byte(r))
		}
		return obj.Value{
			D: []obj.Data{{D: sb.String(), T: obj.VStr}},
		}
	default: // Object.
		return obj.Value{
			D: []obj.Data{{D: fmt.Sprint(args[0].V), T: obj.VStr}},
		}
	}
}

// Append source values to destination array.
func Append(tk obj.Token, args []obj.Var) obj.Value {
	src := args[0].V
	if !src.Arr {
		fract.Panic(tk, obj.ValuePanic, "\"src\" must be array!")
	}
	for _, d := range args[1].V.D {
		src.D = append(src.D, obj.Data{D: d.D, T: d.T})
	}
	return src
}
