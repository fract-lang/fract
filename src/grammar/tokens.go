package grammar

const (
	// GENERIC TOKENS

	// TokenSharp Sharp.
	TokenSharp string = "#"
	// TokenPlus Plus.
	TokenPlus string = "+"
	// TokenMinus Minus.
	TokenMinus string = "-"
	// TokenStar Star.
	TokenStar string = "*"
	// TokenPercent Percent.
	TokenPercent string = "%"
	// TokenSlash Slash.
	TokenSlash string = "/"
	// TokenReverseSlash Rever slash.
	TokenReverseSlash string = "\\"
	// TokenEquals Equals.
	TokenEquals string = "="
	// TokenQuestion Question mark.
	TokenQuestion string = "?"
	// TokenVerticalBar Vertical bar.
	TokenVerticalBar string = "|"
	// TokenGreat Greater then.
	TokenGreat string = ">"
	// TokenLess Less then.
	TokenLess string = "<"
	// TokenSemicolon Semicolon.
	TokenSemicolon string = ";"
	// TokenColon Colon.
	TokenColon string = ":"
	// TokenComma Comma.
	TokenComma string = ","
	// TokenExclamation Exclamation.
	TokenExclamation string = "!"
	// TokenAmper Amper.
	TokenAmper string = "&"
	// TokenAt At.
	TokenAt string = "@"
	// TokenDot Dot.
	TokenDot string = "."
	// TokenLParenthes Left parentheses.
	TokenLParenthes string = "("
	// TokenRParenthes Right parentheses.
	TokenRParenthes string = ")"
	// TokenCaret Caret.
	TokenCaret string = "^"

	// MULTICHAR OPERATORS

	// SeperatorSub Sub element seperator.
	SeperatorSub string = "::"

	// IntegerDivision Integer divide.
	IntegerDivision string = "//"

	// IntegerDivideWithBigger Integer division with bigger.
	IntegerDivideWithBigger string = "\\\\"

	// KEYWORDS

	// KwImport Import packages.
	KwImport string = "use"
	// KwFunction Function define.
	KwFunction string = "fn"
	// KwDelete Delete variable.
	KwDelete string = "del"
	// KwVariable Variable define.
	KwVariable string = "var"
	// KwBlockFinal Block terminator.
	KwBlockFinal string = "end"
	// KwReturn Returns.
	KwReturn string = "ret"
	// KwForLoop For loop.
	KwForLoop string = "for"
	// KwWhileLoop While loop.
	KwWhileLoop string = "while"
	// KwIf If condition.
	KwIf string = "if"
	// KwElseIf Else-If alternate.
	KwElseIf string = "elif"
	// KwElse Else.
	KwElse string = "else"

	// DATA TYPES

	// DtByte byte.
	DtByte string = "uint8"
	// DtSignedByte sbyte.
	DtSignedByte string = "int8"
	// Dt16BitInteger short.
	Dt16BitInteger string = "int16"
	// Dt32BitInteger int.
	Dt32BitInteger string = "int32"
	// Dt64BitInteger long.
	Dt64BitInteger string = "int64"
	// DtUnsigned16BitInteger ushort.
	DtUnsigned16BitInteger string = "uint16"
	// DtUnsigned32BitInteger uint.
	DtUnsigned32BitInteger string = "uint32"
	// DtUnsigned64BitInteger ulong.
	DtUnsigned64BitInteger string = "uint64"
	// DtBoolean boolean.
	DtBoolean string = "bool"
	// DtFloat float.
	DtFloat string = "float"
	// DtDouble double.
	DtDouble string = "double"
)
