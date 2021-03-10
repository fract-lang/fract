/*
	Insert Function.
*/

package vector

// Insert Insert value by position.
// pos Position to insert.
// value Value to insert.
func (v *Vector) Insert(pos int, value ...interface{}) {
	v.Vals = append(v.Vals[:pos], append(value, v.Vals[pos:]...)...)
}
