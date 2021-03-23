/*
	Do Function.
*/

package except

// Do Do call block.
func (block Block) Do() {
	if block.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				block.Catch(r)
			}
		}()
	}
	block.Try()
}
