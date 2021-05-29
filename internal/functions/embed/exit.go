package embed

import (
	"os"
	"strconv"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// Exit from application with code.
func Exit(f objects.Function, parameters []objects.Variable) {
	code := parameters[0].Value
	if code.Array {
		fract.Error(f.Tokens[0][0], "Array is not a valid value!")
	} else if code.Content[0].Type != fract.VALInteger {
		fract.Error(f.Tokens[0][0], "Exit code is only be integer!")
	}
	exit_code, _ := strconv.ParseInt(code.Content[0].Data, 10, 64)
	os.Exit(int(exit_code))
}
