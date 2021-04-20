/*
	ProcessArrayIndex Function.
*/

package parser

// ProcessArrayIndex Process array index by length.
// length Length of array.
// index Index to process.
func ProcessArrayIndex(length, index int) int {
	if index >= 0 {
		if index >= length {
			return -1
		}
		return index
	}

	index = length + index
	if index < 0 || index >= length {
		return -1
	}
	return index
}
