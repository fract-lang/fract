/*
	New Function.
*/

package interpreter

import (
	"../lexer"
	"../utilities/vector"
)

// New Create new instance of Parser.
// path Path of destination file.
// type Type of file.
func New(path string, _type int) *Interpreter {
	preter := new(Interpreter)
	preter.lexer = lexer.New(ReadyFile(path))
	preter.vars = vector.New()
	preter.funcs = vector.New()
	preter.Type = _type
	preter.tokens = vector.New()
	return preter
}
