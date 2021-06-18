package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fract-lang/fract/pkg/str"
)

// Input returns input from command-line.
func Input(message string) string {
	fmt.Print(message)
	//! Don't use fmt.Scanln
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// PrintMapAsTable print map to cli screen as table.
func PrintMapAsTable(dict map[string]string) {
	maxlen := 0
	for key := range dict {
		if maxlen < len(key) {
			maxlen = len(key)
		}
	}
	maxlen += 5
	for key := range dict {
		fmt.Println(key + " " + str.Whitespace(maxlen-len(key)) + dict[key])
	}
}
