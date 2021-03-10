/*
	New Function.
*/

package vector

// New Create new instance.
// values Base values.
func New(values ...interface{}) *Vector {
	vector := new(Vector)
	vector.Vals = append(vector.Vals, values...)
	return vector
}
