package value

import (
	"fmt"
)

const (
	Single uint8 = 0
)

// Val intance.
type Val struct {
	D []Data // Data.
	T uint8  // Type.
}

func (v *Val) String() string {
	if v.D == nil {
		return ""
	}
	return v.D[0].Format()
}

func (v *Val) Print() bool {
	if v.D == nil {
		return false
	}
	fmt.Print(v.String())
	return true
}
