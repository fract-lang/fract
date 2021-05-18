package fs

import "os"

func ExistsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
