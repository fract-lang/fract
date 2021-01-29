package cli

import (
	"bufio"
	"fmt"
	"os"
)

// Returns string whitespace by count.
// count Count of whitespace.
func getws(count int) string {
	var str string = ""
	for counter := 1; counter <= count; counter++ {
		str += " "
	}
	return str
}

// PrintMapAsTable Print map to cli screen as table.
// dict Map to print.
func PrintMapAsTable(dict map[string]string) {
	var maxlen int = 0
	for key := range dict {
		if maxlen < len(key) {
			maxlen = len(key)
		}
	}
	maxlen += 5
	for key := range dict {
		fmt.Println(key + " " + getws(maxlen-len(key)) +
			dict[key])
	}
}

// Input Returns input from command-line.
// message Input message.
func Input(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	text, _ := reader.ReadString('\n')
	return text
}
