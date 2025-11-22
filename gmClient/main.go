package main

import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/gameLogic"
	GM "home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"

)

func main(){
	gameState , _ ,err := gameLogic.NewGamestateWithRandomMap(10, 10, 0.2, []float64{0.5, 0.2, 0.2, 0.1})
	if err != nil{
		fmt.Println("Error creating game state:", err)
		return
	}
	currentMap := gameState.CurrentMap
	err = serialization.SaveMapToFile(currentMap, currentMap.Name)
	if err != nil{
		fmt.Println("Error creating game state:", err)
		return
	}

	gm := GM.NewGM("GameMaster", "The overseer of the game world")
	
	fmt.Println("Welcome,", gm.Name)
	fmt.Println("You are the Game Master.")
	fmt.Println("Type 'map' to view the map, or 'quit' to exit.")

	commands := io.GetInput()
	for {
		command := commands[0]
		switch command {
			case "map":
				currentMap.PrintMapDebug()
			case "quit":
				fmt.Println("Quitting game.")
				return
			default:
				fmt.Println("Unknown command:", command)
		}
		commands = io.GetInput()
	}

}