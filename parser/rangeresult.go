package parser

import "../objects"

// RangeResult Result instance of LexRange function.
type RangeResult struct {
	// Range is found.
	Found bool
	// Tokens of range.
	Range []objects.Token
	// Index of replace index of original list.
	Index int32
}
