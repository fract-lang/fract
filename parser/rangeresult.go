package parser

import "container/list"

// RangeResult Result instance of LexRange function.
type RangeResult struct {
	// Range is found.
	found bool
	// Tokens of range.
	tokens list.List
	// Index of replace index of original list.
	index int
}
