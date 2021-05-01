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
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
)

var preter *interpreter.Interpreter

func interpret() {
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
				preter.Lexer.File.Lines = append(preter.Lexer.File.Lines, interpreter.ReadyLines([]string{input})...)
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
				preter.Lexer.File.Lines = append(preter.Lexer.File.Lines, interpreter.ReadyLines([]string{input})...)
				goto reTokenize
			}

			preter.Tokens = append(preter.Tokens, cacheTokens)
		}

		// Change blocks.
		count := 0
		for _, tokens := range preter.Tokens {
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
			preter.Lexer.File.Lines = append(preter.Lexer.File.Lines, interpreter.ReadyLines([]string{input})...)
			goto reTokenize
		}

		preter.Interpret()
	}
}

func catch(e *objects.Exception) {
	if e.Message != "" {
		fmt.Println("Fract is panicked, sorry this is a problem with Fract!")
		fmt.Println(e.Message)
	}
}

func processCommand(ns, cmd string) {
	switch ns {
	case "help":
		ModuleHelp.Process(cmd)
	case "version":
		ModuleVersion.Process(cmd)
	case "make":
		ModuleMake.Process(cmd)
	default:
		if ModuleMake.Check(ns) {
			ModuleMake.Process(ns + cmd)
		} else {
			fmt.Println("There is no such command!")
		}
	}
}

func init() {
	// Not started with arguments.
	if len(os.Args) < 2 {
		return
	}

	var sb strings.Builder
	for _, current := range os.Args[1:] { sb.WriteString(" " + current) }
	os.Args[0] = sb.String()[1:]

	processCommand(commands.GetNamespace(os.Args[0]), commands.RemoveNamespace(os.Args[0]))

	os.Exit(0)
}

func main() {
	fmt.Println("Fract " + fract.FractVersion + " (c) MIT License.\n" + "Developed by Fract Developer Team.\n")

	if info, err := os.Stat("std"); err != nil || !info.IsDir() {
		fmt.Println("Standard library not found!")
		return
	}

	fract.LiveInterpret = true

	preter = interpreter.NewStdin(".")
	preter.ApplyEmbedFunctions()

	block := new(except.Block)
	block.Try = interpret
	block.Catch = catch

	for {
		block.Do()
	}
}
