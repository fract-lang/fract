package parser

// ProcessArrayIndex process array index by length.
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
