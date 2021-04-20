/*
	GetNamespace Function.
*/

package commands

import "strings"

// GetNamespace Get namespace of command.
// command Command.
func GetNamespace(command string) string {
	position := strings.Index(command, " ")
	if position == -1 {
		return command
	}
	return command[0:position]
}
