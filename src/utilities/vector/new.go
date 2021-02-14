/*
	New Function.
*/

package vector

// New Create new instance.
// values Base values.
func New(values ...interface{}) *Vector {
	vector := new(Vector)
	vector.Vals = make([]interface{}, len(values))
	copy(vector.Vals, values)
	return vector
}
