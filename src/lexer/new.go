/*
	New Function.
*/

package lexer

import "../objects"

// New Create new instance.
func New(file objects.CodeFile) *Lexer {
	lexer := new(Lexer)
	lexer.File = &file
	lexer.Line = 1
	return lexer
}
