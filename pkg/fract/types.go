package fract

const (
	TypeNone   uint8 = 0
	TypeIgnore uint8 = 1

	TypeComment             uint8 = 10
	TypeOperator            uint8 = 11
	TypeValue               uint8 = 12
	TypeBrace               uint8 = 13
	TypeVariable            uint8 = 14
	TypeName                uint8 = 15
	TypeDelete              uint8 = 16
	TypeComma               uint8 = 17
	TypeBlockEnd            uint8 = 18
	TypeIf                  uint8 = 19
	TypeElseIf              uint8 = 20
	TypeElse                uint8 = 21
	TypeStatementTerminator uint8 = 22
	TypeLoop                uint8 = 23
	TypeIn                  uint8 = 24
	TypeBreak               uint8 = 25
	TypeContinue            uint8 = 26
	TypeFunction            uint8 = 27
	TypeReturn              uint8 = 28
	TypeProtected           uint8 = 29
	TypeTry                 uint8 = 30
	TypeCatch               uint8 = 31
	TypeImport              uint8 = 32
	TypeParams              uint8 = 33
	TypeMacro               uint8 = 34
	TypeDefer               uint8 = 35

	LOOPBreak    uint8 = 1
	LOOPContinue uint8 = 2

	FUNCReturn uint8 = 3

	VALInteger uint8 = 0
	VALFloat   uint8 = 1
	VALString  uint8 = 2
	VALBoolean uint8 = 3
)
