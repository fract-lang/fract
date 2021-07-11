package fract

const (
	Ver         = "0.0.1"  // Version of Fract.
	Ext         = ".fract" // File extension of Fract.
	FloatFormat = "%g"     // Float format.

	None                uint8 = 0
	Ignore              uint8 = 1
	Comment             uint8 = 10
	Operator            uint8 = 11
	Value               uint8 = 12
	Brace               uint8 = 13
	Var                 uint8 = 14
	Name                uint8 = 15
	Delete              uint8 = 16
	Comma               uint8 = 17
	If                  uint8 = 18
	Else                uint8 = 19
	StatementTerminator uint8 = 20
	Loop                uint8 = 21
	In                  uint8 = 22
	Break               uint8 = 23
	Continue            uint8 = 24
	Func                uint8 = 25
	Ret                 uint8 = 26
	Protected           uint8 = 27
	Try                 uint8 = 28
	Catch               uint8 = 29
	Import              uint8 = 30
	Params              uint8 = 31
	Macro               uint8 = 32
	Defer               uint8 = 33
	Go                  uint8 = 34

	LOOPBreak    uint8 = 1
	LOOPContinue uint8 = 2
	FUNCReturn   uint8 = 3
)

var (
	TryCount      int // Try-Catch count.
	ExecPath      string
	InteractiveSh bool // Interactive shell mode.
)
