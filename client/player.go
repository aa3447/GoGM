package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main(){
	commands  := getInput()
	for {
		command := commands[0]
		//args := commands[1:]
		switch command {
		case "move":
			//handleMove(args)
		case "action":
			//handleAction(args)
		case "quit":
			fmt.Println("Quitting game.")
			return
		default:
			fmt.Println("Unknown command:", command)
		}
	}
}

func getInput() []string{
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