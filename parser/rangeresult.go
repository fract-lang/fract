package parser

import "container/list"

// RangeResult Result instance of LexRange function.
type RangeResult struct {
	// Range is found.
	Found bool
	// Tokens of range.
	Tokens list.List
	// Index of replace index of original list.
	Index int
}
