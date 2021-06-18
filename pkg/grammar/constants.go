package grammar

// TODO: Optimize here.

const (
	LogicalAnd                 = "&&"
	LogicalOr                  = "||"
	RangeCommentOpen           = "#>"
	RangeCommentClose          = "<#"
	IntegerDivision            = "//"
	IntegerDivideWithBigger    = "\\\\"
	Exponentiation             = "**"
	Equals                     = "=="
	NotEquals                  = "<>"
	GreaterEquals              = ">="
	LessEquals                 = "<="
	Params                     = "..."
	LeftBinaryShift            = "<<"
	RightBinaryShift           = ">>"
	AdditionAssignment         = "+="
	SubtractionAssignment      = "-="
	MultiplicationAssignment   = "*="
	ExponentiationAssignment   = "**="
	DivisionAssignment         = "/="
	ModulusAssignment          = "%="
	LeftBinaryShiftAssignment  = "<<="
	RightBinaryShiftAssignment = ">>="
	InclusiveOrAssignment      = "|="
	XOrAssignment              = "^="
	AndAssignment              = "&="

	KwNaN          = "NaN"       // Not a number.
	KwImport       = "open"      // Import packages.
	KwFunction     = "func"      // Function declaration.
	KwDelete       = "del"       // Delete define(s).
	KwVariable     = "var"       // Variable define.
	KwConstant     = "const"     // Constant variable declaration.
	KwProtected    = "protected" // Protect from manuel memory deletion.
	KwReturn       = "ret"       // Return value.
	KwBlockEnd     = "end"       // End of block.
	KwForWhileLoop = "for"       // For and while loop.
	KwIn           = "in"        // In.
	KwBreak        = "break"     // Break loop.
	KwContinue     = "continue"  // Continue loop.
	KwIf           = "if"        // If condition.
	KwElseIf       = "elif"      // Else-If alternate.
	KwElse         = "else"      // Else.
	KwTrue         = "true"      // Boolean true(1) value.
	KwFalse        = "false"     // Boolean false(0) value.
	KwTry          = "try"       // Try declare.
	KwCatch        = "catch"     // Catch declare.
	KwDefer        = "defer"     // Defer.
	KwMut          = "mut"       // Mutable variable declaration.
)
