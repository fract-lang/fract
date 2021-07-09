package obj

import (
	"fmt"
	"os"
)

const (
	PlainPanic        = "Panic"
	NamePanic         = "NamePanic"
	MemoryPanic       = "MemoryPanic"
	SyntaxPanic       = "SyntaxPanic"
	ValuePanic        = "ValuePanic"
	OutOfRangePanic   = "OutOfRangePanic"
	ArithmeticPanic   = "ArithmeticPanic"
	DivideByZeroPanic = "DivideByZeroPanic"
)

type Panic struct {
	M string // Message.
	T string // Type.
}

func (p Panic) String() string { return p.M }

func (p Panic) Panic() { fmt.Println(p); os.Exit(1) }
