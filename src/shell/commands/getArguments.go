/*
	GetArguments Function.
*/

package commands

import (
	"regexp"

	"../../utils/vector"
)

// GetArguments Get arguments of command.
// command Command.
func GetArguments(command string) vector.Vector {
	var args vector.Vector
	pattern := regexp.MustCompile("(^|\\s+)-\\w+(?=($|\\s+))")
	for arg := range pattern.FindAllString(command, -1) {
		args.Vals = append(args.Vals, arg)
	}
	return args
}
