package cli

import (
	"bufio"
	"fmt"
	"os"
)

// Input returns input from command-line.
func Input(message string) string {
	fmt.Print(message)
	//! Don't use fmt.Scanln
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
