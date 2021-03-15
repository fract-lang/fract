/*
	isKeywordToken Function.
*/

package lexer

import (
	"regexp"
)

// isKeywordToken Returns true if statement is keyword compatible token, false if not.
// ln Line.
// kw Target keyword.
func isKeywordToken(ln, kw string) bool {
	return regexp.MustCompile(
		"^" + kw + `(\s+|$|[[:punct:]])`).MatchString(ln)
}
