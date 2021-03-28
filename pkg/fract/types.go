/*
	TYPES OF FRACT
*/

package fract

const (
	// TypeNone NA
	TypeNone = -1
	// TypeIgnore Ignore.
	TypeIgnore = 1
	// TypeEntryFile Entry file.
	TypeEntryFile = 1000
	// TypeImportedFile Imported file.
	TypeImportedFile = 1001

	// TypeComment Comment.
	TypeComment = 1100
	// TypeOperator Operator.
	TypeOperator = 1101
	// TypePrint Print.
	TypePrint = 1102
	// TypeValue Value.
	TypeValue = 1103
	// TypeBrace Bracket.
	TypeBrace = 1104
	// TypeVariable Variable define.
	TypeVariable = 1105
	// TypeName Name type.
	TypeName = 1106
	// TypeDelete Delete from memory.
	TypeDelete = 1107
	// TypeComma Comma.
	TypeComma = 1108
	// TypeBooleanTrue Boolean true.
	TypeBooleanTrue = 119
	// TypeBooleanFalse Boolean false.
	TypeBooleanFalse = 1110
	// TypeBlockEnd End of block.
	TypeBlockEnd = 1111
	// TypeIf If.
	TypeIf = 1112
	// TypeElseIf Else if.
	TypeElseIf = 1113
	// TypeElse Else.
	TypeElse = 1114
	// TypeStatementTerminator Statement terminator.
	TypeStatementTerminator = 1115
	// TypeLoop Loop.
	TypeLoop = 1116
	// TypeIn In.
	TypeIn = 1117
	// TypeBreak Break loop.
	TypeBreak = 1118
	// TypeContinue Continue loop.
	TypeContinue = 1119
	// TypeFunction Function declare.
	TypeFunction = 1120
	// TypeReturn Return.
	TypeReturn = 1121
	// TypeProtected Protected.
	TypeProtected = 1122
	// TypeTry Try.
	TypeTry = 1123
	// TypeCatch Catch.
	TypeCatch = 1124
	// TypeImport Import.
	TypeImport = 1125

	// LOOPBreak Break loop.
	LOOPBreak = 0
	// LOOPContinue Continue loop.
	LOOPContinue = 1

	// FUNCReturn Return value.
	FUNCReturn = 2

	// VALInteger Integer value.
	VALInteger = 0
	// VALFloat Float value.
	VALFloat = 1
	// VALString String value.
	VALString = 2
	// Boolean Boolean value.
	VALBoolean = 3
)
