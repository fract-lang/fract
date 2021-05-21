package fs

import "io/ioutil"

func ReadAllText(path string) string {
	content, err := ioutil.ReadFile(path)
	if err == nil {
		return string(content)
	}
	return ""
}
