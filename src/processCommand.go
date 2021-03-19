/*
	processCommand Function.
*/

package main

import (
	"fmt"

	ModuleHelp "github.com/fract-lang/fract/src/shell/modules/help"
	ModuleMake "github.com/fract-lang/fract/src/shell/modules/make"
	ModuleVersion "github.com/fract-lang/fract/src/shell/modules/version"
)

func processCommand(ns, cmd string) {
	if ns == "help" {
		ModuleHelp.Process(cmd)
	} else if ns == "version" {
		ModuleVersion.Process(cmd)
	} else if ns == "make" {
		ModuleMake.Process(cmd)
	} else if ModuleMake.Check(ns) {
		ModuleMake.Process(ns + cmd)
	} else {
		fmt.Println("There is no such command!")
	}
}
