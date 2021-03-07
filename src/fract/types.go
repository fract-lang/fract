/*
	TYPES OF FRACT
*/

package fract

const (
	// TypeNone NA
	TypeNone = -1
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
	// TypeBlock Block start type.
	TypeBlock = 1111
	// TypeBlockEnd End of block type.
	TypeBlockEnd = 1112
	// TypeIf If declare.
	TypeIf = 1113
	// TypeElseIf Else if declare.
	TypeElseIf = 1114
	// TypeElse Else declare.
	TypeElse = 1115
	// TypeStatementTerminator Statement terminator.
	TypeStatementTerminator = 1116
	// TypeLoop Loop type.
	TypeLoop = 1117
	// TypeIn In type.
	TypeIn = 1118
	// TypeBreak Break loop type.
	TypeBreak = 1119
	// TypeContinue Continue loop type.
	TypeContinue = 1120
	// TypeExit Exit type.
	TypeExit = 1121
	// TypeFunction Function declare type.
	TypeFunction = 1122
	// TypeReturn Return type.
	TypeReturn = 1123

	// LOOPBreak Break loop.
	LOOPBreak = 0
	// LOOPContinue Continue loop.
	LOOPContinue = 1

	// FUNCReturn Return value.
	FUNCReturn = 2
)
