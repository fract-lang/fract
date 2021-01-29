package fs

import (
	"io/ioutil"
	"os"
	"strings"
)

// Rename Rename folder or file by path.
// path Path to rename.
// newName New name of path.
func Rename(path string, newName string) {
	stat, _ := os.Stat(path)
	newPath := path[0:len(path)-len(stat.Name())] + newName
	os.Rename(path, newPath)
}

// ExistsPath Returns true if path is exits, returns false if not.
// path Path to check.
func ExistsPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

// ExistsFile Returns true if file is exits, returns false if not.
// path Path to check.
func ExistsFile(path string) bool {
	_, err := ioutil.ReadFile(path)
	if err == nil {
		return true
	}
	return false
}

// ExistsDirectory Returns true if folder is exits, returns false if not.
// path Path to check.
func ExistsDirectory(path string) bool {
	_, err := ioutil.ReadDir(path)
	if err == nil {
		return true
	}
	return false
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
	return strings.Split(ReadAllText(path), '\n')
}

// WriteText Append to content.
// path Destination path.
// content Content to write.
func WriteText(path string, content string) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	file.WriteString(content)
	file.Close()
}

// WriteAllText Remove contetn and write this content.
// path Destination path.
// content Content to write.
func WriteAllText(path string, content string) {
	file, _ := os.Create(path)
	file.WriteString(content)
	file.Close()
}

// ReadAllBytes Read bytes from content by path.
// path Path to read.
func ReadAllBytes(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err == nil {
		return content
	}
	return nil
}
