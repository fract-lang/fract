package vector

// Vector dynamic list.
// Values Values of vector.
// Len Count of element.
type Vector struct {
	Vals []interface{}
}

// New Create new instance.
// values Base values.
func New(values ...interface{}) *Vector {
	vector := new(Vector)
	vector.Vals = values
	return vector
}

// Len Returns length of vector.
func (v *Vector) Len() int {
	return len(v.Vals)
}

// Append Append value.
// value Value to append.
func (v *Vector) Append(value ...interface{}) {
	v.Vals = append(v.Vals, value...)
}

// RemoveFirst Remove first element.
func (v *Vector) RemoveFirst() {
	v.Vals = v.Vals[1:]
}

// RemoveLast Remove last element.
func (v *Vector) RemoveLast() {
	v.Vals = v.Vals[:len(v.Vals)-1]
}

// Insert Insert value by position.
// pos Position to insert.
// value Value to insert.
func (v *Vector) Insert(pos int, value ...interface{}) {
	len := len(value)
	v.Vals = append(v.Vals, len)
	copy(v.Vals[pos+len:], v.Vals[pos:])
	for counter := 0; counter < len; counter++ {
		v.Vals[pos] = value[counter]
		pos++
	}
}

// At Get element by index.
// pos Index of element.
func (v *Vector) At(pos int) interface{} {
	return v.Vals[pos]
}

// AtCheck Get element by index and check.
// pos Index of element.
func (v *Vector) AtCheck(pos int) (interface{}, bool) {
	if pos < 0 || pos >= len(v.Vals) {
		return nil, false
	}
	return v.Vals[pos], true
}

// Set Set value by index.
// pos Index to set.
func (v *Vector) Set(pos int, value interface{}) {
	v.Vals[pos] = value
}

// RemoveRange Remove range.
// pos Start position of removing.
// len Count of removing elements.
func (v *Vector) RemoveRange(pos int, len int) {
	v.Vals = append(v.Vals[:pos], v.Vals[pos+len:]...)
}

// Remove Remove first matched element.
// value Value instance to remove.
func (v *Vector) Remove(value interface{}) bool {
	for index := 0; index < len(v.Vals); index++ {
		if v.Vals[index] == value {
			v.Vals = append(v.Vals[:index], v.Vals[index+1:]...)
			return true
		}
	}
	return false
}

// RemoveAll Remove all matched elements.
// value Value instance to remove.
func (v *Vector) RemoveAll(value interface{}) bool {
	result := false
	for v.Remove(value) {
		result = true
	}
	return result
}

// Find Find element.
// pos Start position to search.
// value Value to search.
func (v *Vector) Find(pos int, value interface{}) int {
	for ; pos < len(v.Vals); pos++ {
		if v.Vals[pos] == value {
			return pos
		}
	}
	return -1
}

// Exist Check exist element.
// value Value to check.
func (v *Vector) Exist(value interface{}) bool {
	for index := 0; index < len(v.Vals); index++ {
		if v.Vals[index] == value {
			return true
		}
	}
	return false
}

// Reverse Reverse elements.
func (v *Vector) Reverse() {
	len := len(v.Vals)
	for index := 0; index < len/2; index++ {
		cache := v.Vals[index]
		v.Vals[index] = v.Vals[len-index-1]
		v.Vals[len-index-1] = cache
	}
}

// Clear Remove all elements.
func (v *Vector) Clear() {
	v.Vals = make([]interface{}, 0)
}

// Any Any value is exist?
func (v *Vector) Any() bool {
	return len(v.Vals) > 0
}

// Sublist Get range.
// pos Start position to take.
// len Count of taken elements.
func (v *Vector) Sublist(pos int, len int) Vector {
	sub := New()
	for counter := 1; counter <= len; counter++ {
		sub.Append(v.Vals[pos])
		pos++
	}
	return *sub
}

// First Returns first element.
func (v *Vector) First() interface{} {
	return v.Vals[0]
}

// Last Returns last element.
func (v *Vector) Last() interface{} {
	return v.Vals[len(v.Vals)-1]
}

// Join Concatenate vectors.
// vector Vector to concatenate.
func (v *Vector) Join(vector Vector) {
	v.Vals = append(v.Vals, vector.Vals...)
}

// Equals Equals all elements and order to vector?
// vector Vector to check.
func (v *Vector) Equals(vector Vector) bool {
	xlen := len(v.Vals)
	ylen := len(vector.Vals)

	if xlen != ylen {
		return false
	}

	for index := 0; index < xlen; index++ {
		if v.Vals[index] != vector.Vals[index] {
			return false
		}
	}

	return true
}

// Unique Remove copies of repeated elements.
func (v *Vector) Unique() {
	vtr := New()
	for index := 0; index < len(v.Vals); index++ {
		current := v.Vals[index]
		if vtr.Exist(current) {
			continue
		}
		vtr.Append(current)
	}
	v.Vals = vtr.Vals
}
