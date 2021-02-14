/*
  RemoveArguments Function.
*/

package commands

import "regexp"

// RemoveArguments Remove arguments from command.
// command Command.
func RemoveArguments(command string) string {
	pattern := regexp.MustCompile("(^|\\s+)-\\w+(?=($|\\s+))")
	return pattern.ReplaceAllString(command, "")
}
