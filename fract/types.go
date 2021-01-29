package fract

const (
	// TypeNone NA
	TypeNone int = -1
	// TypeEntryFile Entry file.
	TypeEntryFile int = 99
	// TypeImportedFile Imported file.
	TypeImportedFile int = 100
	// TypeComment Comment.
	TypeComment int = 999
	// TypeFunction Function.
	TypeFunction int = 1000
	// TypeEquals Equals.
	TypeEquals int = 1001
	// TypeLet Let.
	TypeLet int = 1002
	// TypeName Name.
	TypeName int = 1003
	// TypeDottedName Dotted name.
	TypeDottedName int = 1004
	// TypeValueSetter Value setter.
	TypeValueSetter int = 1005
	// TypeValue Value.
	TypeValue int = 1006
	// TypeDataType Data type.
	TypeDataType int = 1007
	// TypeEndType End type.
	TypeEndType int = 1008
	// TypeReturn Return.
	TypeReturn int = 1009
	// TypeImport Import.
	TypeImport int = 1010
	// TypeStdImport Standard library import.
	TypeStdImport int = 1011
	// TypeIf If condition.
	TypeIf int = 1012
	// TypeElseIf Else-If condition.
	TypeElseIf int = 1013
	// TypeElse Else condition.
	TypeElse int = 1014
	// TypeFor For loop.
	TypeFor int = 1015
	// TypeWhile While loop.
	TypeWhile int = 1016
	// TypeDelete Delete.
	TypeDelete int = 1017
	// TypeShort 16Bit integer.
	TypeShort int = 1018
	// TypeInt 32Bit integer.
	TypeInt int = 1019
	// TypeLong 64Bit integer.
	TypeLong int = 1020
	// TypeUShort Unsigned 16Bit integer.
	TypeUShort int = 1021
	// TypeUInt Unsigned 32Bit integer.
	TypeUInt int = 1022
	// TypeULong Unsigned 64Bit integer.
	TypeULong int = 1023
	// TypeFloat 16Bit float.
	TypeFloat int = 1024
	// TypeDouble 32Bit float.
	TypeDouble int = 1025
	// TypeBoolean Boolean.
	TypeBoolean int = 1026
	// TypeByte Unsigned 8Bit integer.
	TypeByte int = 1027
	// TypeSByte 8Bit integer.
	TypeSByte int = 1028
	// TypeOperator Operator.
	TypeOperator int = 1029
	// TypeOpenParenthes Open parenthes.
	TypeOpenParenthes int = 1030
	// TypeCloseParenthes Close parenthes.
	TypeCloseParenthes int = 1031
)
