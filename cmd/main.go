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
	"strings"

	"github.com/fract-lang/fract/internal/interpreter"
	"github.com/fract-lang/fract/internal/shell/commands"
	ModuleHelp "github.com/fract-lang/fract/internal/shell/modules/help"
	ModuleMake "github.com/fract-lang/fract/internal/shell/modules/make"
	ModuleVersion "github.com/fract-lang/fract/internal/shell/modules/version"
	"github.com/fract-lang/fract/pkg/cli"
	"github.com/fract-lang/fract/pkg/except"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/parser"
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

func init() {
	// Not started with arguments.
	if len(os.Args) < 2 {
		return
	}

	var sb strings.Builder
	os.Args = os.Args[1:]
	for _, current := range os.Args {
		sb.WriteString(" " + current)
	}
	os.Args[0] = sb.String()[1:]

	processCommand(commands.GetNamespace(os.Args[0]),
		commands.RemoveNamespace(os.Args[0]))

	os.Exit(0)
}

func main() {
	fmt.Println(
		"Fract " + fract.FractVersion + " (c) MIT License.\n" +
			"Developed by Fract Developer Team.\n")

	preter := interpreter.NewStdin()
	preter.ApplyEmbedFunctions()

repeat:
	except.Block{
		Try: func() {
			for {
				input := cli.Input(">> ")
				preter.Lexer.File.Lines = interpreter.ReadyLines([]string{input})
			reTokenize:
				preter.Tokens = nil
			reTokenizeUnNil:
				preter.Lexer.Finished = false
				preter.Lexer.Line = 1
				preter.Lexer.Column = 1

				/* Tokenize all lines. */
				for !preter.Lexer.Finished {
					cacheTokens := preter.Lexer.Next()

					// Check multiline comment.
					if preter.Lexer.RangeComment {
						input := cli.Input(" | ")
						preter.Lexer.File.Lines = append(preter.Lexer.File.Lines,
							interpreter.ReadyLines([]string{input})...)
						goto reTokenizeUnNil
					}

					// cacheTokens are empty?
					if cacheTokens == nil {
						continue
					}

					// Check parentheses.
					if preter.Lexer.BraceCount > 0 ||
						preter.Lexer.BracketCount > 0 ||
						preter.Lexer.ParenthesCount > 0 {
						input := cli.Input(" | ")
						preter.Lexer.File.Lines = append(preter.Lexer.File.Lines,
							interpreter.ReadyLines([]string{input})...)
						goto reTokenize
					}

					preter.Tokens = append(preter.Tokens, cacheTokens)
				}

				// Change blocks.
				count := 0
				for index := range preter.Tokens {
					tokens := preter.Tokens[index]
					if first := tokens[0]; first.Type == fract.TypeBlockEnd {
						count--
						if count < 0 {
							fract.Error(first, "The extra block end defined!")
						}
					} else if parser.IsBlockStatement(tokens) {
						count++
					}
				}

				if count > 0 { // Check blocks.
					input := cli.Input(" | ")
					preter.Lexer.File.Lines = append(preter.Lexer.File.Lines,
						interpreter.ReadyLines([]string{input})...)
					goto reTokenize
				}

				preter.Interpret()
			}
		},
		/*Catch: func(e except.Exception) {
			if e != "" {
				fmt.Println("Fract is panicked, sorry this is a problem with Fract!")
				fmt.Println(e)
			}
		},*/
	}.Do()
	goto repeat
}
