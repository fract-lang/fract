package array

// Remove Remove range from array.
func Remove(array *[]interface{}, start int, count int) {
	copy((*array)[start:], (*array)[start+count:])
	(*array)[len(*array)-1] = nil
	*array = (*array)[:len(*array)-1]
}
