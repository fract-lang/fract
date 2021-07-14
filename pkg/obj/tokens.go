package obj

// Token instance.
type Token struct {
	F   *File
	V   string
	T   uint8
	Ln  int
	Col int
}

type Tokens []Token

// Sub slice.
func (t Tokens) Sub(pos, len int) *Tokens {
	if len == 0 {
		return nil
	}
	t = append([]Token{}, t[pos:pos+len]...)
	return &t
}

// Remove range.
func (t *Tokens) Rem(pos, len int) {
	if len > 0 {
		*t = append((*t)[:pos], (*t)[pos+len:]...)
	}
}

// Insert at.
func (t *Tokens) Ins(pos int, vals ...Token) {
	*t = append((*t)[:pos], append(vals, (*t)[pos:]...)...)
}
