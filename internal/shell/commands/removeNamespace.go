/*
	RemoveNamespace Function.
*/

package commands

import "strings"

// RemoveNamespace Remove namespace from command.
// command Command.
func RemoveNamespace(command string) string {
	position := strings.Index(command, " ")
	if position == -1 {
		return ""
	}
	return command[position+1:]
}
