/*
	Do Function.
*/

package except

import "github.com/fract-lang/fract/pkg/objects"

// Do Do call block.
func (block *Block) Do() {
	defer func() {
		if r := recover(); r != nil {
			block.Exception = &objects.Exception{
				Message: r.(error).Error(),
			}

			if block.Catch != nil {
				block.Catch(block.Exception)
			}
		}
	}()
	block.Try()
}
