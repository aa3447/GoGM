package main

import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/gameLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
)

func main(){
	gameState , entranceLocation, err := gameLogic.NewGamestateWithRandomMap(10, 10, 0.2, []float64{0.5, 0.2, 0.2, 0.1})
	if err != nil{
		fmt.Println("Error creating game state:", err)
		return
	}
	err = serialization.SaveMapToFile(gameState.CurrentMap, gameState.CurrentMap.Name)
	if err != nil{
		fmt.Println("Error creating game state:", err)
		return
	}
	player := player.NewPlayer("Hero", "The brave adventurer", "Warrior")
	player.SetLocation(entranceLocation[0], entranceLocation[1])
	
	fmt.Println("Welcome,", player.Name)
	fmt.Println("You find yourself at the entrance of a mysterious location.")
	fmt.Println("Type 'move <direction>' to move (north, south, east, west), 'map' to view the map, or 'quit' to exit.")
	
	commands := io.GetInput()
	for {
		command := commands[0]
		args := commands[1:]
		switch command {
			case "move":
				handleMove(gameState,player,args)
			case "action":
				//handleAction(args)
			case "map":
				gameState.CurrentMap.PrintMapWithPlayer(player)
			case "quit":
				fmt.Println("Quitting game.")
				return
			default:
				fmt.Println("Unknown command:", command)
		}
		commands = io.GetInput()
	}
}


func handleMove(gs *gameLogic.Gamestate , player *player.Player ,args []string){
	if len(args) < 1{
		fmt.Println("Move where?")
		return
	}
	direction := args[0]
	fmt.Println("Moving", direction)

	switch direction{
		case "north", "up":
			gs.MovePlayer(player,-1, 0)
		case "south", "down":
			gs.MovePlayer(player,1, 0)
		case "east", "right":
			gs.MovePlayer(player,0, 1)
		case "west", "left":
			gs.MovePlayer(player,0, -1)
		default:
			fmt.Println("Unknown direction:", direction)
	}
}