/*
	ARGUMENT FUNCTIONS
*/

package commands

import (
	"regexp"

	"../../utilities/vector"
)

// GetArguments Get arguments of command.
// command Command.
func GetArguments(command string) vector.Vector {
	var args vector.Vector
	pattern := regexp.MustCompile("(^|\\s+)-\\w+(?=($|\\s+))")
	for arg := range pattern.FindAllString(command, -1) {
		args.Append(arg)
	}
	return args
}

// RemoveArguments Remove arguments from command.
// command Command.
func RemoveArguments(command string) string {
	pattern := regexp.MustCompile("(^|\\s+)-\\w+(?=($|\\s+))")
	return pattern.ReplaceAllString(command, "")
}
