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
	// TypeBrace Bracket type.
	TypeBrace = 1104
	// TypeVariable Variable define type.
	TypeVariable = 1105
	// TypeName Name type.
	TypeName = 1106
	// TypeDelete Delete from memory type.
	TypeDelete = 1107
	// TypeComma Comma type.
	TypeComma = 1108
	// TypeBooleanTrue Boolean true type.
	TypeBooleanTrue = 119
	// TypeBooleanFalse Boolean false type.
	TypeBooleanFalse = 1110
	// TypeBlockEnd End of block type.
	TypeBlockEnd = 1111
	// TypeIf If declare.
	TypeIf = 1112
	// TypeElseIf Else if declare.
	TypeElseIf = 1113
	// TypeElse Else declare.
	TypeElse = 1114
	// TypeStatementTerminator Statement terminator.
	TypeStatementTerminator = 1115
	// TypeLoop Loop type.
	TypeLoop = 1116
	// TypeIn In type.
	TypeIn = 1117
	// TypeBreak Break loop type.
	TypeBreak = 1118
	// TypeContinue Continue loop type.
	TypeContinue = 1119
	// TypeExit Exit type.
	TypeExit = 1120
	// TypeFunction Function declare type.
	TypeFunction = 1121
	// TypeReturn Return type.
	TypeReturn = 1122
	// TypeProtected Protected type.
	TypeProtected = 1123
	// TypeTry Try type.
	TypeTry = 1124
	// TypeCatch Catch type.
	TypeCatch = 1125

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
