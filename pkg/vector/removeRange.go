/*
	RemoveRange Function.
*/

package vector

// RemoveRange Remove range.
// pos Start position of removing.
// len Count of removing elements.
func (v *Vector) RemoveRange(pos, len int) {
	v.Vals = append(v.Vals[:pos], v.Vals[pos+len:]...)
}
