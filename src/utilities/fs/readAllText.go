/*
	ReadAllText Function.
*/

package fs

import (
	"io/ioutil"
)

// ReadAllText Read all text from file by path.
// path Path to read.
func ReadAllText(path string) string {
	content, err := ioutil.ReadFile(path)
	if err == nil {
		return string(content)
	}
	return ""
}
