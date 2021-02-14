/*
	Insert Function.
*/

package vector

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
