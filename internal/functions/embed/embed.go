package embed

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
func Exit(f obj.Func, args []obj.Var) {
	c := args[0].Val
	if c.Arr {
		fract.Error(f.Tks[0][0], "Array is not a valid value!")
	} else if c.D[0].T != obj.VInt {
		fract.Error(f.Tks[0][0], "Exit code is only be integer!")
	}
	ec, _ := strconv.ParseInt(c.D[0].String(), 10, 64)
	os.Exit(int(ec))
}

// Float convert object to float.
func Float(f obj.Func, parameters []obj.Var) obj.Value {
	return obj.Value{D: []obj.Data{
		{D: fmt.Sprintf(fract.FloatFormat, arithmetic.Value(parameters[0].Val.D[0].String())), T: obj.VFloat},
	}}
}

// Input returns input from command-line.
func Input(f obj.Func, args []obj.Var) obj.Value {
	args[0].Val.Print()
	//! Don't use fmt.Scanln
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	return obj.Value{
		D: []obj.Data{{D: s.Text(), T: obj.VStr}},
	}
}

// Int convert object to integer.
func Int(f obj.Func, args []obj.Var) obj.Value {
	switch args[1].Val.D[0].D { // Cast type.
	case "strcode":
		var v obj.Value
		for _, byt := range []byte(args[0].Val.D[0].String()) {
			v.D = append(v.D, obj.Data{D: fmt.Sprint(byt), T: obj.VInt})
		}
		v.Arr = len(v.D) > 1
		return v
	default: // Object.
		return obj.Value{
			D: []obj.Data{{D: fmt.Sprint(int(arithmetic.Value(args[0].Val.D[0].String()))), T: obj.VInt}},
		}
	}
}

// Len returns length of object.
func Len(f obj.Func, args []obj.Var) obj.Value {
	arg := args[0].Val
	if arg.Arr {
		return obj.Value{D: []obj.Data{{D: fmt.Sprint(len(arg.D))}}}
	} else if arg.D[0].T == obj.VStr {
		return obj.Value{D: []obj.Data{{D: fmt.Sprint(len(arg.D[0].String()))}}}
	}
	return obj.Value{D: []obj.Data{{D: "0"}}}
}

// Make array by size.
func Make(f obj.Func, args []obj.Var) obj.Value {
	sz := args[0].Val
	if sz.Arr {
		fract.Error(f.Tks[0][0], "Array is not a valid value!")
	} else if sz.D[0].T != obj.VInt {
		fract.Error(f.Tks[0][0], "Exit code is only be integer!")
	}
	szv, _ := strconv.Atoi(sz.D[0].String())
	if szv < 0 {
		fract.Error(f.Tks[0][0], "Size should be minimum zero!")
	}
	v := obj.Value{Arr: true}
	if szv > 0 {
		var index int
		for ; index < szv; index++ {
			v.D = append(v.D, obj.Data{D: "0"})
		}
	} else {
		v.D = []obj.Data{}
	}
	return v
}

// Print values to cli.
func Print(f obj.Func, args []obj.Var) {
	if args[0].Val.D == nil {
		fract.Error(f.Tks[0][0], "Invalid value!")
	}
	args[0].Val.Print()
	args[1].Val.Print()
}

// Range returns array by parameters.
func Range(f obj.Func, args []obj.Var) obj.Value {
	start := args[0].Val
	to := args[1].Val
	step := args[2].Val
	if start.Arr {
		fract.Error(f.Tks[0][0], "\"start\" argument should be numeric!")
	} else if to.Arr {
		fract.Error(f.Tks[0][0], "\"to\" argument should be numeric!")
	} else if step.Arr {
		fract.Error(f.Tks[0][0], "\"step\" argument should be numeric!")
	}
	if start.D[0].T != obj.VInt && start.D[0].T != obj.VFloat || to.D[0].T != obj.VInt &&
		to.D[0].T != obj.VFloat || step.D[0].T != obj.VInt && step.D[0].T != obj.VFloat {
		fract.Error(f.Tks[0][0], "Values should be integer or float!")
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
func String(f obj.Func, args []obj.Var) obj.Value {
	switch args[1].Val.D[0].D {
	case "parse":
		str := ""
		if value := args[0].Val; value.Arr {
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
			str = args[0].Val.D[0].String()
		}
		return obj.Value{
			D: []obj.Data{{D: str, T: obj.VStr}},
		}
	case "bytecode":
		v := args[0].Val
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
			D: []obj.Data{{D: fmt.Sprint(args[0].Val), T: obj.VStr}},
		}
	}
}

// Append source values to destination array.
func Append(f obj.Func, args []obj.Var) obj.Value {
	src := args[0].Val
	if !src.Arr {
		fract.Error(f.Tks[0][0], "\"src\" must be array!")
	}
	for _, d := range args[1].Val.D {
		src.D = append(src.D, obj.Data{D: d.D, T: d.T})
	}
	return src
}
