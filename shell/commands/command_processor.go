package commands

import (
	"container/list"
	"regexp"
	"strings"
)

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

// GetArguments Get arguments of command.
// command Command.
func GetArguments(command string) list.List {
	var args list.List
	pattern := regexp.MustCompile("(^|\\s+)-\\w+(?=($|\\s+))")
	for arg := range pattern.FindAllString(command, -1) {
		args.PushBack(arg)
	}
	return args
}

// RemoveArguments Remove arguments from command.
// command Command.
func RemoveArguments(command string) string {
	pattern := regexp.MustCompile("(^|\\s+)-\\w+(?=($|\\s+))")
	return pattern.ReplaceAllString(command, "")
}
