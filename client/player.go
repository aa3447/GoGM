package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/gameLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
)

func PlayerLoop(){
	commands := getInput()
	gameState , err := gameLogic.NewGamestateWithRandomMap(10, 10, 0.2, []float64{0.5, 0.2, 0.2, 0.1})
	if err != nil{
		fmt.Println("Error creating game state:", err)
		return
	}
	for {
		command := commands[0]
		args := commands[1:]
		switch command {
			case "move":
				handleMove(gameState,args)
			case "action":
				//handleAction(args)
			case "map":
				playerY, playerX := gameState.GetPlayerPosition()
				mapLogic.PrintMapWithPlayer(gameState.CurrentMap, playerY, playerX)
			case "quit":
				fmt.Println("Quitting game.")
				return
			default:
				fmt.Println("Unknown command:", command)
		}
		commands = getInput()
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

func handleMove(gs *gameLogic.Gamestate ,args []string){
	if len(args) < 1{
		fmt.Println("Move where?")
		return
	}
	direction := args[0]
	fmt.Println("Moving", direction)

	switch direction{
		case "north", "up":
			gs.MovePlayer(-1, 0)
		case "south", "down":
			gs.MovePlayer(1, 0)
		case "east", "right":
			gs.MovePlayer(0, 1)
		case "west", "left":
			gs.MovePlayer(0, -1)
		default:
			fmt.Println("Unknown direction:", direction)
	}
}