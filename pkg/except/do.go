package except

import "github.com/fract-lang/fract/pkg/objects"

func (block *Block) catch() {
	if r := recover(); r != nil {
		block.Exception = &objects.Exception{
			Message: r.(error).Error(),
		}

		if block.Catch != nil {
			block.Catch(block.Exception)
		}
	}
}

// Do execute block.
func (block *Block) Do() {
	defer block.catch()
	block.Try()
}
