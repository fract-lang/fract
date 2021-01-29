package cli

import (
	"bufio"
	"fmt"
	"os"
)

// Input Returns input from command-line.
// message Input message.
func Input(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	text, _ := reader.ReadString('\n')
	return text
}
