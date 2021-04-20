/*
	Do Function.
*/

package except

import "github.com/fract-lang/fract/pkg/objects"

// Do Do call block.
func (block Block) Do() {
	if block.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				block.Catch(objects.Exception{
					Message: r.(error).Error(),
				})
			}
		}()
	}
	block.Try()
}
