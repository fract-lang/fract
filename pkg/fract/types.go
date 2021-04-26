/*
	TYPES OF FRACT
*/

package fract

const (
	// TypeNone NA
	TypeNone int16 = -1
	// TypeIgnore Ignore.
	TypeIgnore int16 = 1

	// TypeComment Comment.
	TypeComment int16 = 1100
	// TypeOperator Operator.
	TypeOperator int16 = 1101
	// TypePrint Print.
	TypePrint int16 = 1102
	// TypeValue Value.
	TypeValue int16 = 1103
	// TypeBrace Bracket.
	TypeBrace int16 = 1104
	// TypeVariable Variable define.
	TypeVariable int16 = 1105
	// TypeName Name type.
	TypeName int16 = 1106
	// TypeDelete Delete from memory.
	TypeDelete int16 = 1107
	// TypeComma Comma.
	TypeComma int16 = 1108
	// TypeBooleanTrue Boolean true.
	TypeBooleanTrue int16 = 119
	// TypeBooleanFalse Boolean false.
	TypeBooleanFalse int16 = 1110
	// TypeBlockEnd End of block.
	TypeBlockEnd int16 = 1111
	// TypeIf If.
	TypeIf int16 = 1112
	// TypeElseIf Else if.
	TypeElseIf int16 = 1113
	// TypeElse Else.
	TypeElse int16 = 1114
	// TypeStatementTerminator Statement terminator.
	TypeStatementTerminator int16 = 1115
	// TypeLoop Loop.
	TypeLoop int16 = 1116
	// TypeIn In.
	TypeIn int16 = 1117
	// TypeBreak Break loop.
	TypeBreak int16 = 1118
	// TypeContinue Continue loop.
	TypeContinue int16 = 1119
	// TypeFunction Function declare.
	TypeFunction int16 = 1120
	// TypeReturn Return.
	TypeReturn int16 = 1121
	// TypeProtected Protected.
	TypeProtected int16 = 1122
	// TypeTry Try.
	TypeTry int16 = 1123
	// TypeCatch Catch.
	TypeCatch int16 = 1124
	// TypeImport Import.
	TypeImport int16 = 1125
	// TypeParams Params.
	TypeParams int16 = 1126

	// LOOPBreak Break loop.
	LOOPBreak int16 = 0
	// LOOPContinue Continue loop.
	LOOPContinue int16 = 1

	// FUNCReturn Return value.
	FUNCReturn int16 = 2

	// VALInteger Integer value.
	VALInteger int16 = 0
	// VALFloat Float value.
	VALFloat int16 = 1
	// VALString String value.
	VALString int16 = 2
	// Boolean Boolean value.
	VALBoolean int16 = 3
)
