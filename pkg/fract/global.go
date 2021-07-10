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
	End                 uint8 = 18
	If                  uint8 = 19
	ElseIf              uint8 = 20
	Else                uint8 = 21
	StatementTerminator uint8 = 22
	Loop                uint8 = 23
	In                  uint8 = 24
	Break               uint8 = 25
	Continue            uint8 = 26
	Func                uint8 = 27
	Ret                 uint8 = 28
	Protected           uint8 = 29
	Try                 uint8 = 30
	Catch               uint8 = 31
	Import              uint8 = 32
	Params              uint8 = 33
	Macro               uint8 = 34
	Defer               uint8 = 35
	Go                  uint8 = 36

	LOOPBreak    uint8 = 1
	LOOPContinue uint8 = 2
	FUNCReturn   uint8 = 3
)

var (
	TryCount      int // Try-Catch count.
	ExecPath      string
	InteractiveSh bool // Interactive shell mode.
)
