/*
	ERROR FUNCTIONS
*/

package lexer

import (
	"fmt"
	"os"
)

// Error Exit with error.
// message Message of error.
func (l *Lexer) Error(message string) {
	fmt.Printf("LEXER ERROR\nMessage: %s\nLINE: %d\nCOLUMN: %d",
		message, l.Line, l.Column)
	os.Exit(1)
}
