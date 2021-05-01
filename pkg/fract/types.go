package fract

const (
	TypeNone   uint8 = 0
	TypeIgnore uint8 = 1

	TypeComment             uint8 = 100
	TypeOperator            uint8 = 101
	TypePrint               uint8 = 102 // TODO: Check this type is unecessary.
	TypeValue               uint8 = 103
	TypeBrace               uint8 = 104
	TypeVariable            uint8 = 105 // Variable define.
	TypeName                uint8 = 106
	TypeDelete              uint8 = 107
	TypeComma               uint8 = 108
	TypeBooleanTrue         uint8 = 109
	TypeBooleanFalse        uint8 = 110
	TypeBlockEnd            uint8 = 111
	TypeIf                  uint8 = 112
	TypeElseIf              uint8 = 113
	TypeElse                uint8 = 114
	TypeStatementTerminator uint8 = 115
	TypeLoop                uint8 = 116
	TypeIn                  uint8 = 117
	TypeBreak               uint8 = 118
	TypeContinue            uint8 = 119
	TypeFunction            uint8 = 120 // Function define.
	TypeReturn              uint8 = 121
	TypeProtected           uint8 = 122
	TypeTry                 uint8 = 123
	TypeCatch               uint8 = 124
	TypeImport              uint8 = 125
	TypeParams              uint8 = 126

	LOOPBreak    uint8 = 1
	LOOPContinue uint8 = 2

	FUNCReturn uint8 = 3

	VALInteger uint8 = 0
	VALFloat   uint8 = 1
	VALString  uint8 = 2
	VALBoolean uint8 = 3
)
