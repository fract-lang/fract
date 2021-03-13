/*
	Sublist Function.
*/

package vector

// Sublist Get range.
// pos Start position to take.
// length Count of taken elements.
func (v *Vector) Sublist(pos int, length int) *Vector {
	if length == 0 {
		return &Vector{}
	}
	return &Vector{Vals: v.Vals[pos : pos+length]}
}
