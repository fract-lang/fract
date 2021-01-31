// Copyright (c) 2021 Fract
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"os"

	"./shell/commands"

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

func main() {
	// not started with arguments.
	if len(os.Args) < 2 {
		os.Exit(0)
	}

	command := ""
	skipped := false
	for arg := range os.Args {
		if !skipped {
			skipped = true
			continue
		}
		if command == "" {
			command += os.Args[arg]
			continue
		}
		command += " " + os.Args[arg]
	}

	processCommand(commands.GetNamespace(command),
		commands.RemoveNamespace(command))
}
