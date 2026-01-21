package io

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetInput reads a line of input from the standard input and returns it as a slice of strings.
func GetInput() []string{
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanned := scanner.Scan()
	if !scanned {
		return nil
	}
	line := scanner.Text()
	line = strings.TrimSpace(line)
	line = strings.ToLower(line)
	return strings.Fields(line)
}

// ConfirmationWrapper wraps a function that requires confirmation before proceeding.
func ConfirmationWrapper(wrapped func() error, prompt string) bool{
	if err := wrapped(); err != nil {
		return false
	}
	for {
		fmt.Printf("%s (y/n): ", prompt)
		input := GetInput()
		if len(input) == 0 {
			continue
		}
		response := input[0]
		switch response {
			case "y", "Y":
				return true
			case "n", "N":
				if err := wrapped(); err != nil {
					return false
				}
			default:
				fmt.Println("Please enter 'y' or 'n'.")
			}
	}
}