/*
	Sublist Function.
*/

package vector

// Sublist Get range.
// pos Start position to take.
// length Count of taken elements.
func (v *Vector) Sublist(pos, length int) *Vector {
	if length == 0 {
		return &Vector{}
	}
	return &Vector{
		Vals: append(make([]interface{}, 0), v.Vals[pos:pos+length]...),
	}
}
