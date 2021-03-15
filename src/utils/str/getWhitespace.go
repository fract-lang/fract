package str

// Returns string whitespace by count.
// count Count of whitespace.
func GetWhitespace(count int) string {
	str := ""
	for counter := 1; counter <= count; counter++ {
		str += " "
	}
	return str
}
