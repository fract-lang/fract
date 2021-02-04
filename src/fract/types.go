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
	// TypeDataType Data type type.
	TypeDataType int = 1107

	// VTInteger Integer value type.
	VTInteger int = 0
	// VTFloat Float value type.
	VTFloat int = 1
)
