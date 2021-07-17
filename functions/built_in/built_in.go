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
	if c.T != value.Int {
		fract.Panic(tk, obj.ValuePanic, "Exit code is only be integer!")
	}
	ec, _ := strconv.ParseInt(c.String(), 10, 64)
	os.Exit(int(ec))
}

// Float convert object to float.
func Float(parameters []obj.Var) value.Val {
	return value.Val{
		D: fmt.Sprintf(fract.FloatFormat, value.Conv(parameters[0].V.String())),
		T: value.Float,
	}
}

// Input returns input from command-line.
func Input(args []obj.Var) value.Val {
	args[0].V.Print()
	//! Don't use fmt.Scanln
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	return value.Val{D: s.Text(), T: value.Str}
}

// Int convert object to integer.
func Int(args []obj.Var) value.Val {
	switch args[1].V.D { // Cast type.
	case "strcode":
		var v []value.Val
		for _, byt := range []byte(args[0].V.String()) {
			v = append(v, value.Val{D: fmt.Sprint(byt), T: value.Int})
		}
		return value.Val{D: v, T: value.Array}
	default: // Object.
		return value.Val{
			D: fmt.Sprint(int(value.Conv(args[0].V.String()))),
			T: value.Int,
		}
	}
}

// Len returns length of object.
func Len(args []obj.Var) value.Val {
	return value.Val{D: fmt.Sprint(args[0].V.Len()), T: value.Int}
}

// Calloc array by size.
func Calloc(tk obj.Token, args []obj.Var) value.Val {
	sz := args[0].V
	if sz.T != value.Int {
		fract.Panic(tk, obj.ValuePanic, "Size is only be integer!")
	}
	szv, _ := strconv.Atoi(sz.String())
	if szv < 0 {
		fract.Panic(tk, obj.ValuePanic, "Size should be minimum zero!")
	}
	v := value.Val{T: value.Array}
	if szv > 0 {
		var index int
		var data []value.Val
		for ; index < szv; index++ {
			data = append(data, value.Val{D: "0", T: value.Int})
		}
		v.D = data
	} else {
		v.D = []value.Val{}
	}
	return v
}

// Realloc array by size.
func Realloc(tk obj.Token, args []obj.Var) value.Val {
	if args[0].V.T != value.Array {
		fract.Panic(tk, obj.ValuePanic, "Value is must be array!")
	}
	szv, _ := strconv.Atoi(args[1].V.String())
	if szv < 0 {
		fract.Panic(tk, obj.ValuePanic, "Size should be minimum zero!")
	}
	var (
		data []value.Val
		b    = args[0].V.D.([]value.Val)
		v    = value.Val{T: value.Array}
		c    = 0
	)
	if len(b) <= szv {
		data = b
		c = len(b)
	} else {
		v.D = b[:szv]
		return v
	}
	for ; c <= szv; c++ {
		data = append(data, value.Val{D: "0", T: value.Int})
	}
	v.D = data
	return v
}

// Print values to cli.
func Print(tk obj.Token, args []obj.Var) {
	if args[0].V.D == nil {
		fract.Panic(tk, obj.ValuePanic, "Value is not printable!")
	}
	for _, d := range args[0].V.D.([]value.Val) {
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
	if start.T != value.Int && start.T != value.Float {
		fract.Panic(tk, obj.ValuePanic, "\"start\" argument should be numeric!")
	} else if to.T != value.Int && to.T != value.Float {
		fract.Panic(tk, obj.ValuePanic, "\"to\" argument should be numeric!")
	} else if step.T != value.Int && step.T != value.Float {
		fract.Panic(tk, obj.ValuePanic, "\"step\" argument should be numeric!")
	}
	if start.T != value.Int && start.T != value.Float || to.T != value.Int &&
		to.T != value.Float || step.T != value.Int && step.T != value.Float {
		fract.Panic(tk, obj.ValuePanic, "Values should be integer or float!")
	}
	startV, _ := strconv.ParseFloat(start.String(), 64)
	toV, _ := strconv.ParseFloat(to.String(), 64)
	stepV, _ := strconv.ParseFloat(step.String(), 64)
	if stepV <= 0 {
		return value.Val{T: value.Array}
	}
	var t uint8
	if start.T == value.Float || to.T == value.Float || step.T == value.Float {
		t = value.Float
	}
	var data []value.Val
	if startV <= toV {
		for ; startV <= toV; startV += stepV {
			d := value.Val{D: fmt.Sprintf(fract.FloatFormat, startV), T: t}
			d.D = d.Format()
			data = append(data, d)
		}
	} else {
		for ; startV >= toV; startV -= stepV {
			d := value.Val{D: fmt.Sprintf(fract.FloatFormat, startV), T: t}
			d.D = d.Format()
			data = append(data, d)
		}
	}
	return value.Val{D: data, T: value.Array}
}

// String convert object to string.
func String(args []obj.Var) value.Val {
	switch args[1].V.D {
	case "parse":
		str := ""
		if val := args[0].V; val.T == value.Array {
			data := val.D.([]value.Val)
			if len(data) == 0 {
				str = "[]"
			} else {
				var sb strings.Builder
				sb.WriteByte('[')
				for _, data := range data {
					sb.WriteString(data.String() + " ")
				}
				str = sb.String()[:sb.Len()-1] + "]"
			}
		} else {
			str = args[0].V.String()
		}
		return value.Val{D: str, T: value.Str}
	case "bytecode":
		v := args[0].V
		var sb strings.Builder
		for _, d := range v.D.([]value.Val) {
			if d.T != value.Int {
				sb.WriteByte(' ')
			}
			r, _ := strconv.ParseInt(d.String(), 10, 32)
			sb.WriteByte(byte(r))
		}
		return value.Val{D: sb.String(), T: value.Str}
	default: // Object.
		return value.Val{D: fmt.Sprint(args[0].V), T: value.Str}
	}
}

// Append source values to destination array.
func Append(tk obj.Token, args []obj.Var) value.Val {
	src := args[0].V
	if src.T != value.Array {
		fract.Panic(tk, obj.ValuePanic, "\"src\" must be array!")
	}
	data := args[1].V.D.([]value.Val)
	for _, d := range data {
		data = append(data, value.Val{D: d.D, T: d.T})
	}
	src.D = data
	return src
}
