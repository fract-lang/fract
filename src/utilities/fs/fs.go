package fs

import (
	"io/ioutil"
	"os"
	"strings"
)

// ExistFile Returns true if file is exist, returns false if not.
// path Path to check.
func ExistFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// ReadAllText Read all text from file by path.
// path Path to read.
func ReadAllText(path string) string {
	content, err := ioutil.ReadFile(path)
	if err == nil {
		return string(content)
	}
	return ""
}

// ReadAllLines Read all lines from file by path.
// path Path to read.
func ReadAllLines(path string) []string {
	return strings.Split(ReadAllText(path), "\n")
}

// OpenFile Open file by path.
// path Path to open.
func OpenFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return file
	}
	return nil
}
