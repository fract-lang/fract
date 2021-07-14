package main

import (
	"testing"

	"github.com/fract-lang/fract/parser"
)

func BenchmarkInterpret(b *testing.B) {
	p := parser.New("../test.fract")
	p.AddBuiltInFuncs()
	for i := 0; i < b.N; i++ {
		p.Interpret()
	}
}
