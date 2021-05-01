package fract

const (
	TypeNone   uint8 = 0
	TypeIgnore uint8 = 1

	TypeComment             uint8 = 10
	TypeOperator            uint8 = 11
	TypeValue               uint8 = 12
	TypeBrace               uint8 = 13
	TypeVariable            uint8 = 14 // Variable define.
	TypeName                uint8 = 15
	TypeDelete              uint8 = 16
	TypeComma               uint8 = 17
	TypeBooleanTrue         uint8 = 18
	TypeBooleanFalse        uint8 = 19
	TypeBlockEnd            uint8 = 20
	TypeIf                  uint8 = 21
	TypeElseIf              uint8 = 22
	TypeElse                uint8 = 23
	TypeStatementTerminator uint8 = 24
	TypeLoop                uint8 = 25
	TypeIn                  uint8 = 26
	TypeBreak               uint8 = 27
	TypeContinue            uint8 = 28
	TypeFunction            uint8 = 29 // Function define.
	TypeReturn              uint8 = 30
	TypeProtected           uint8 = 31
	TypeTry                 uint8 = 32
	TypeCatch               uint8 = 33
	TypeImport              uint8 = 34
	TypeParams              uint8 = 35
	TypeMacro               uint8 = 36

	LOOPBreak    uint8 = 1
	LOOPContinue uint8 = 2

	FUNCReturn uint8 = 3

	VALInteger uint8 = 0
	VALFloat   uint8 = 1
	VALString  uint8 = 2
	VALBoolean uint8 = 3
)
