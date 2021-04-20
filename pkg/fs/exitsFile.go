/*
	ExitsFile Function.
*/

package fs

import (
	"os"
)

// ExistFile Returns true if file is exist, returns false if not.
// path Path to check.
func ExistFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
