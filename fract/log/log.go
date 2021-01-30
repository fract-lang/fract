package log

import (
	"../../utilities/list"
)

var (
	// GlobalLog Global log instance.
	GlobalLog *Log = New()
)

// Log Log instance of Fract.
type Log struct {
	/* PUBLIC */

	// Logs.
	Values list.List
}

// New Create new instance.
func New() *Log {
	return new(Log)
}

// Add Add log.
// log Log to add.
func (l *Log) Add(log MessageLog) {
	l.Values.Append(log)
}

// Error Add error log.
// message Message of error.
// line Line of error.
func (l *Log) Error(message string, line int) {
	l.Values.Append(MessageLog{Message: message, Line: line})
}

// Clear Clear all logs.
func (l *Log) Clear() {
	l.Values.Clear()
}
