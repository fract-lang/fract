package commands

import "strings"

// GetNamespace returns namespace of command.
func GetNamespace(command string) string {
	position := strings.Index(command, " ")
	if position == -1 {
		return command
	}
	return command[0:position]
}
