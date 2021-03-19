/*
	OpenFile Function.
*/

package fs

import (
	"os"
)

// OpenFile Open file by path.
// path Path to open.
func OpenFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return file
	}
	return nil
}
