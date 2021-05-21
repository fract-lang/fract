package commands

import "strings"

// RemoveNamespace remove namespace from command.
func RemoveNamespace(command string) string {
	position := strings.Index(command, " ")
	if position == -1 {
		return ""
	}
	return command[position+1:]
}
