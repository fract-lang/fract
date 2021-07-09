package obj

// Code block instance.
type Block struct {
	Try   func()
	Catch func(Panic)
	P     Panic
}

func (b *Block) catch() {
	if r := recover(); r != nil {
		b.P = r.(Panic)
		if b.Catch != nil {
			b.Catch(b.P)
		}
	}
}

// Do execute block.
func (b *Block) Do() {
	defer b.catch()
	b.Try()
}
