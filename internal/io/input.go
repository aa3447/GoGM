package io

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetInput() []string{
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanned := scanner.Scan()
	if !scanned {
		return nil
	}
	line := scanner.Text()
	line = strings.TrimSpace(line)
	return strings.Fields(line)
}