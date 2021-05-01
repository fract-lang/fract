package fract

const (
	TypeNone   uint8 = 0
	TypeIgnore uint8 = 1

	TypeComment             uint8 = 100
	TypeOperator            uint8 = 101
	TypeValue               uint8 = 102
	TypeBrace               uint8 = 103
	TypeVariable            uint8 = 104 // Variable define.
	TypeName                uint8 = 105
	TypeDelete              uint8 = 106
	TypeComma               uint8 = 107
	TypeBooleanTrue         uint8 = 108
	TypeBooleanFalse        uint8 = 109
	TypeBlockEnd            uint8 = 110
	TypeIf                  uint8 = 111
	TypeElseIf              uint8 = 112
	TypeElse                uint8 = 113
	TypeStatementTerminator uint8 = 114
	TypeLoop                uint8 = 115
	TypeIn                  uint8 = 116
	TypeBreak               uint8 = 117
	TypeContinue            uint8 = 118
	TypeFunction            uint8 = 119 // Function define.
	TypeReturn              uint8 = 120
	TypeProtected           uint8 = 121
	TypeTry                 uint8 = 122
	TypeCatch               uint8 = 123
	TypeImport              uint8 = 124
	TypeParams              uint8 = 125

	LOOPBreak    uint8 = 1
	LOOPContinue uint8 = 2

	FUNCReturn uint8 = 3

	VALInteger uint8 = 0
	VALFloat   uint8 = 1
	VALString  uint8 = 2
	VALBoolean uint8 = 3
)
