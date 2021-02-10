/*
	TYPES OF FRACT
*/

package fract

const (
	// TypeNone NA
	TypeNone int = -1
	// TypeEntryFile Entry file.
	TypeEntryFile int = 1000
	// TypeImportedFile Imported file.
	TypeImportedFile int = 1001

	// TypeComment Comment.
	TypeComment int = 1100
	// TypeOperator Operator.
	TypeOperator int = 1101
	// TypePrint Print.
	TypePrint int = 1102
	// TypeValue Value.
	TypeValue int = 1103
	// TypeBrace Bracket type.
	TypeBrace int = 1104
	// TypeVariable Variable define type.
	TypeVariable int = 1105
	// TypeName Name type.
	TypeName int = 1106
	// TypeDataType Datatype type.
	TypeDataType int = 1107
	// TypeDelete Delete from memory type.
	TypeDelete int = 1108
	// TypeComma Comma type.
	TypeComma int = 1109
	// TypeBooleanTrue Boolean true type.
	TypeBooleanTrue int = 1110
	// TypeBooleanFalse Boolean false type.
	TypeBooleanFalse int = 1111
	// TypeBlock Block start type.
	TypeBlock int = 1112
	// TypeBlockEnd End of block type.
	TypeBlockEnd int = 1113
	// TypeIf If declare.
	TypeIf int = 1114
	// TypeElseIf Else if/else declare.
	TypeElseIf int = 1115
	// TypeStatementTerminator Statement terminator.
	TypeStatementTerminator int = 1116

	// VTInteger Integer value type.
	VTInteger int = 0
	// VTIntegerArray Integer array value type.
	VTIntegerArray int = 1
	// VTFloat Float value type.
	VTFloat int = 2
	// VTFloatArray Float array value type.
	VTFloatArray int = 3
)
