/*
	checkName Function.
*/

package interpreter

import (
	fractName "../fract/name"
)

// checkName Check name is exist or not?
// name Name to check.
func (i *Interpreter) checkName(name string) bool {
	return fractName.VarIndexByName(i.vars, name) != -1
}
