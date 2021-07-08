package obj

// Code block instance.
type Block struct {
	Try   func()
	Catch func(Exception)
	E     Exception
}

func (b *Block) catch() {
	if r := recover(); r != nil {
		b.E = Exception{
			Msg: r.(error).Error(),
		}
		if b.Catch != nil {
			b.Catch(b.E)
		}
	}
}

// Do execute block.
func (b *Block) Do() {
	defer b.catch()
	b.Try()
}
