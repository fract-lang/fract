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
