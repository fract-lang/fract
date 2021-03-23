package functions

import (
	"github.com/fract-lang/fract/pkg/cli"
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Input get input from cli.
// f Function.
func Input(f obj.Function) obj.Value {
	fract.PrintValue(f.Parameters[0].Default)
	return obj.Value{
		Content: []string{cli.Input("")},
		Type:    fract.VALString,
	}
}
