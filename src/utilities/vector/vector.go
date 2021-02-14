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

// Append Append value.
// value Value to append.
func (v *Vector) Append(value ...interface{}) {
	v.Vals = append(v.Vals, value...)
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

// RemoveRange Remove range.
// pos Start position of removing.
// len Count of removing elements.
func (v *Vector) RemoveRange(pos int, len int) {
	v.Vals = append(v.Vals[:pos], v.Vals[pos+len:]...)
}

// Sublist Get range.
// pos Start position to take.
// length Count of taken elements.
func (v *Vector) Sublist(pos int, length int) Vector {
	return *New(v.Vals[pos : pos+length]...)
}
