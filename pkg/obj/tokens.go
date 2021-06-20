package obj

type Tokens []Token

func (t Tokens) Sub(pos, len int) *Tokens {
	if len == 0 {
		return nil
	}
	t = append([]Token{}, t[pos:pos+len]...)
	return &t
}

func (t *Tokens) Remove(pos, len int) {
	if len > 0 {
		*t = append((*t)[:pos], (*t)[pos+len:]...)
	}
}

func (t *Tokens) Insert(pos int, vals ...Token) {
	*t = append((*t)[:pos], append(vals, (*t)[pos:]...)...)
}
