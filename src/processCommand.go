/*
	processCommand Function.
*/

package main

import (
	"fmt"

	ModuleExit "./shell/modules/exit"
	ModuleHelp "./shell/modules/help"
	ModuleMake "./shell/modules/make"
	ModuleVersion "./shell/modules/version"
)

func processCommand(ns string, cmd string) {
	if ns == "help" {
		ModuleHelp.Process(cmd)
	} else if ns == "exit" {
		ModuleExit.Process(cmd)
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
