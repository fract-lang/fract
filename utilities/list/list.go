package list

// List dynamic list.
// Values Values of list.
// Len Count of element.
type List struct {
	Vals []interface{}
}

// New Create new instance.
func New() *List {
	return new(List)
}

// Len Returns length of list.
func (l *List) Len() int {
	return len(l.Vals)
}

// Append Append value to list.
// value Value to append.
func (l *List) Append(value interface{}) {
	l.Vals = append(l.Vals, value)
}

// RemoveFirst Remove first element from list.
func (l *List) RemoveFirst() {
	l.Vals = l.Vals[1:]
}

// RemoveLast Remove last element from list.
func (l *List) RemoveLast() {
	l.Vals = l.Vals[:len(l.Vals)-1]
}

// Insert Insert value to list by position.
// pos Position to insert.
// value Value to insert.
func (l *List) Insert(pos int, value interface{}) {
	l.Vals = append(l.Vals, 0)
	copy(l.Vals[pos+1:], l.Vals[pos:])
	l.Vals[pos] = value
}

// At Get element by index.
// pos Index of element.
func (l *List) At(pos int) interface{} {
	return l.Vals[pos]
}

// Set Set value of list by index.
// pos Index to set.
func (l *List) Set(pos int, value interface{}) {
	l.Vals[pos] = value
}

// Remove Remove range from list.
// pos Start position of removing.
// len Count of removing elements.
func (l *List) Remove(pos int, len int) {
	l.Vals = append(l.Vals[:pos], l.Vals[pos+len:]...)
}

// Find Find element.
// pos Start position to search.
// value Value to search.
func (l *List) Find(pos int, value interface{}) int {
	for ; pos < len(l.Vals); pos++ {
		if l.Vals[pos] == value {
			return pos
		}
	}
	return -1
}

// Exist Check exist element in list.
// value Value to check.
func (l *List) Exist(value interface{}) bool {
	for index := 0; index < len(l.Vals); index++ {
		if l.Vals[index] == value {
			return true
		}
	}
	return false
}

// Reverse Reverse list elements.
func (l *List) Reverse() {
	var len int = len(l.Vals)
	for index := 0; index < len/2; index++ {
		var cache interface{} = l.Vals[index]
		l.Vals[index] = l.Vals[len-index-1]
		l.Vals[len-index-1] = cache
	}
}

// Clear Clear all elements from list.
func (l *List) Clear() {
	l.Vals = make([]interface{}, 0)
}

// Any Any value is exist in list?
func (l *List) Any() bool {
	return len(l.Vals) > 0
}

// Sublist Get range from list.
// pos Start position to take.
// len Count of taken elements.
func (l *List) Sublist(pos int, len int) List {
	return List{Vals: l.Vals[pos : len+1]}
}
