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
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fract-lang/fract/internal/parser"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
	"github.com/fract-lang/fract/pkg/str"
)

// Returns ns of command.
func ns(cmd string) string {
	i := strings.Index(cmd, " ")
	if i == -1 {
		return cmd
	}
	return cmd[0:i]
}

// Remove namespace from command.
func remns(cmd string) string {
	i := strings.Index(cmd, " ")
	if i == -1 {
		return ""
	}
	return cmd[i+1:]
}

func input(msg string) string {
	fmt.Print(msg)
	//! Don't use fmt.Scanln
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	return s.Text()
}

var p *parser.Parser = nil

func interpret() {
	for {
		p.L.F.Lns = parser.Lines([]string{input(">> ")})
	reTokenize:
		p.Tks = nil
	reTokenizeUnNil:
		p.L.Fin = false
		p.L.Braces = 0
		p.L.Brackets = 0
		p.L.Parentheses = 0
		p.L.RangeComment = false
		p.L.Ln = 1
		p.L.Col = 1
		/* Tokenize all lines. */
		for !p.L.Fin {
			tks := p.L.Next()
			// Check multiline comment.
			if p.L.RangeComment {
				p.L.F.Lns = append(p.L.F.Lns, parser.Lines([]string{input(" | ")})...)
				goto reTokenizeUnNil
			}
			// cacheTokens are empty?
			if tks == nil {
				continue
			}
			// Check parentheses.
			if p.L.Braces > 0 || p.L.Brackets > 0 || p.L.Parentheses > 0 {
				p.L.F.Lns = append(p.L.F.Lns, parser.Lines([]string{input(" | ")})...)
				goto reTokenize
			}
			p.Tks = append(p.Tks, tks)
		}
		// Change blocks.
		c := 0
		for _, tokens := range p.Tks {
			if first := tokens[0]; first.T == fract.End {
				c--
				if c < 0 {
					fract.IPanic(first, obj.SyntaxPanic, "The extra block end defined!")
				}
			} else if parser.IsBlock(tokens) {
				c++
			}
		}
		if c > 0 { // Check blocks.
			p.L.F.Lns = append(p.L.F.Lns, parser.Lines([]string{input(" | ")})...)
			goto reTokenize
		}
		p.Interpret()
	}
}

func catch(e obj.Panic) {
	if e.M == "" {
		return
	}
	fmt.Println("Fract is panicked, sorry this is a problem with Fract!")
	fmt.Println(e.M)
}

func help(cmd string) {
	if cmd != "" {
		fmt.Println("This module can only be used!")
		return
	}
	d := map[string]string{
		"make":    "Interprete Fract code.",
		"version": "Show version.",
		"help":    "Show help.",
		"exit":    "Exit.",
	}
	mlen := 0
	for k := range d {
		if mlen < len(k) {
			mlen = len(k)
		}
	}
	mlen += 5
	for k := range d {
		fmt.Println(k + " " + str.Whitespace(mlen-len(k)) + d[k])
	}
}

func version(cmd string) {
	if cmd != "" {
		fmt.Println("This module can only be used!")
		return
	}
	fmt.Println("Fract Version [" + fract.Ver + "]")
}

func make(cmd string) {
	if cmd == "" {
		fmt.Println("This module cannot only be used!")
		return
	} else if !strings.HasSuffix(cmd, fract.Ext) {
		cmd += fract.Ext
	}
	if info, err := os.Stat(cmd); err != nil || info.IsDir() {
		fmt.Println("The Fract file is not exists: " + cmd)
		return
	}
	p := parser.New(cmd)
	p.AddBuiltInFuncs()
	(&obj.Block{
		Try: p.Interpret,
		Catch: func(e obj.Panic) {
			os.Exit(0)
		},
	}).Do()
}

func makechk(p string) bool {
	if strings.HasSuffix(p, fract.Ext) {
		return true
	}
	info, err := os.Stat(p + fract.Ext)
	return err == nil && !info.IsDir()
}

func proccmd(ns, cmd string) {
	switch ns {
	case "help":
		help(cmd)
	case "version":
		version(cmd)
	case "make":
		make(cmd)
	default:
		if makechk(ns) {
			make(ns)
		} else {
			fmt.Println("There is no such command!")
		}
	}
}

func init() {
	fract.ExecPath = filepath.Dir(os.Args[0])
	// Check standard library.
	if info, err := os.Stat(path.Join(fract.ExecPath, "std")); err != nil || !info.IsDir() {
		fmt.Println("Standard library not found!")
		input("\nPress enter for exit...")
		os.Exit(1)
	}
	// Not started with arguments.
	if len(os.Args) < 2 {
		return
	}

	defer os.Exit(0)
	var sb strings.Builder
	for _, arg := range os.Args[1:] {
		sb.WriteString(" " + arg)
	}
	os.Args[0] = sb.String()[1:]
	proccmd(ns(os.Args[0]), remns(os.Args[0]))
}

func main() {
	fmt.Println("Fract " + fract.Ver + " (c) MIT License.\n" + "Developed by Fract Developer Team.\n")
	fract.InteractiveSh = true
	p = parser.NewStdin()
	p.AddBuiltInFuncs()
	b := &obj.Block{
		Try:   interpret,
		Catch: catch,
	}
	for {
		b.Do()
	}
}
