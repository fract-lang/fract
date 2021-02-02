/*
	NAMESPACE FUNCTIONS
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

// RemoveNamespace Remove namespace from command.
// command Command.
func RemoveNamespace(command string) string {
	position := strings.Index(command, " ")
	if position == -1 {
		return ""
	}
	return command[position+1:]
}
