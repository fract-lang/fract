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

// ExistPath Returns true if path is exist, returns false if not.
// path Path to check.
func ExistPath(path string) bool {
	_, err := os.Stat(path)
	return err != nil
}

// ExistFile Returns true if file is exist, returns false if not.
// path Path to check.
func ExistFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// ExistDirectory Returns true if folder is exist, returns false if not.
// path Path to check.
func ExistDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
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

// OpenFile Open file by path.
// path Path to open.
func OpenFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return file
	}
	return nil
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
