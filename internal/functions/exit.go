package functions

import (
	"os"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Exit from application with code..
// f Function.
func Exit(f obj.Function) {
	code := f.Parameters[0].Default

	if code.Array {
		fract.Error(f.Tokens[0][0], "Array is not a valid value!")
	} else if code.Type != fract.VALInteger {
		fract.Error(f.Tokens[0][0], "Exit code is only be integer!")
	}

	exit_code, _ := arithmetic.ToInt64(code.Content[0])
	os.Exit(int(exit_code))
}
