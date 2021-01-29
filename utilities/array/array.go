package array

// Remove Remove range from array.
func Remove(array *[]interface{}, start int, count int) {
	copy((*array)[start:], (*array)[start+count:])
	(*array)[len(*array)-1] = nil
	*array = (*array)[:len(*array)-1]
}

// Insert
func Insert(array *[]interface{}, pos int, value interface{}) {
	// Make space for new element.
	*array = append(*array, 0)
	// Insert new element.
	var _type type = value.(type)
	*array = append((*array)[:pos], append([]_type{value}, a[pos:]...)...)
}
