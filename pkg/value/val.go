package value

import (
	"fmt"
)

// Val intance.
type Val struct {
	D   []Data
	Arr bool
}

func (v *Val) String() string {
	if v.D == nil {
		return ""
	}
	if v.Arr {
		return stringArray(v.D)
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
