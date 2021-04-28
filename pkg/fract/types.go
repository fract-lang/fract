/*
	TYPES OF FRACT
*/

package fract

const (
	// TypeNone NA
	TypeNone uint8 = 0
	// TypeIgnore Ignore.
	TypeIgnore uint8 = 1

	// TypeComment Comment.
	TypeComment uint8 = 100
	// TypeOperator Operator.
	TypeOperator uint8 = 101
	// TypePrint Print.
	TypePrint uint8 = 102
	// TypeValue Value.
	TypeValue uint8 = 103
	// TypeBrace Bracket.
	TypeBrace uint8 = 104
	// TypeVariable Variable define.
	TypeVariable uint8 = 105
	// TypeName Name type.
	TypeName uint8 = 106
	// TypeDelete Delete from memory.
	TypeDelete uint8 = 107
	// TypeComma Comma.
	TypeComma uint8 = 108
	// TypeBooleanTrue Boolean true.
	TypeBooleanTrue uint8 = 109
	// TypeBooleanFalse Boolean false.
	TypeBooleanFalse uint8 = 110
	// TypeBlockEnd End of block.
	TypeBlockEnd uint8 = 111
	// TypeIf If.
	TypeIf uint8 = 112
	// TypeElseIf Else if.
	TypeElseIf uint8 = 113
	// TypeElse Else.
	TypeElse uint8 = 114
	// TypeStatementTerminator Statement terminator.
	TypeStatementTerminator uint8 = 115
	// TypeLoop Loop.
	TypeLoop uint8 = 116
	// TypeIn In.
	TypeIn uint8 = 117
	// TypeBreak Break loop.
	TypeBreak uint8 = 118
	// TypeContinue Continue loop.
	TypeContinue uint8 = 119
	// TypeFunction Function declare.
	TypeFunction uint8 = 120
	// TypeReturn Return.
	TypeReturn uint8 = 121
	// TypeProtected Protected.
	TypeProtected uint8 = 122
	// TypeTry Try.
	TypeTry uint8 = 123
	// TypeCatch Catch.
	TypeCatch uint8 = 124
	// TypeImport Import.
	TypeImport uint8 = 125
	// TypeParams Params.
	TypeParams uint8 = 126

	// LOOPBreak Break loop.
	LOOPBreak uint8 = 1
	// LOOPContinue Continue loop.
	LOOPContinue uint8 = 2

	// FUNCReturn Return value.
	FUNCReturn uint8 = 3

	// VALInteger Integer value.
	VALInteger uint8 = 0
	// VALFloat Float value.
	VALFloat uint8 = 1
	// VALString String value.
	VALString uint8 = 2
	// Boolean Boolean value.
	VALBoolean uint8 = 3
)
