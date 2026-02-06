package main

import (
	"fmt"
	"os"
	"errors"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/campaign"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
)


func main(){
	fmt.Println("Load or create a character!")
	var choice string
	var err error
	player := &playerLogic.Player{}
	loadedCampaign := &campaign.Campaign{}
	
	fmt.Println("1. Create new character")
	fmt.Println("2. Load existing character")
	fmt.Scanln(&choice)
	switch choice{
		case "1":
			fmt.Println("Creating new character...")
			player = playerLogic.CreatePlayer()
			// Source - https://stackoverflow.com/a/12518877
			// Posted by Sridhar Ratnakumar, modified by community. See post 'Timeline' for change history
			// Retrieved 2026-02-06, License - CC BY-SA 4.0
			filePath := fmt.Sprintf("./playerClient/player/%s.json", player.Name)
			if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
				serialization.SaveToFile(*player, "player", "player", player.Name)
			}
		case "2":
			fmt.Println("Loading existing character...")
			var name string
			fmt.Print("Enter character name: ")
			fmt.Scanln(&name)
			player, err = serialization.LoadFromJSONFile("player", "player", name, playerLogic.Player{})
			if err != nil{
				fmt.Println("Error loading character:", err)
				return
			}
		default:
			fmt.Println("Invalid choice. Please select 1 or 2.")
	}

	fmt.Println("1. randomly generate campaign")
	fmt.Println("2. Load campaign")
	fmt.Scanln(&choice)
	switch choice{
		case "1":
			fmt.Println("Generating random campaign...")
			gm := playerLogic.GM{Name: "GM", Description: "The Game Master"}
			loadedCampaign = campaign.NewCampaign("Random Campaign", "A randomly generated campaign", gm)
			_, err = loadedCampaign.NewGamestateWithRandomMap(10, 10, 0.2, []float64{0.7, 0.2, 0.1, 0.0}, "Random Map")
			if err != nil{
				fmt.Println("Error generating campaign:", err)
				return
			}
			err = loadedCampaign.SetCurrentMap("Random Map")
			if err != nil{
				fmt.Println("Error setting current map:", err)
				return
			}
		case "2":
			fmt.Println("Loading existing campaign...")
			var campaignName string
			fmt.Print("Enter campaign name: ")
			fmt.Scanln(&campaignName)
			loadedCampaign, err = serialization.LoadFromJSONFile("gm", "campaign", campaignName, campaign.Campaign{})
			if err != nil{
				fmt.Println("Error loading campaign:", err)
				return
			}
		default:
			fmt.Println("Invalid choice. Please select 1 or 2.")
	}
	
	loadedCampaign.Players[player.Name] = player

	fmt.Printf("Starting campaign '%s' with character '%s'...\n", loadedCampaign.Name, player.Name)
	GameLoop(player, loadedCampaign)
}

func GameLoop(player *playerLogic.Player, currentCampaign *campaign.Campaign){
	currentMap := currentCampaign.CurrentMap
	if currentMap == nil{
		fmt.Println("No current map set for the campaign.")
		return
	}

	fmt.Printf("Entering map: %s\n", currentMap.Name)
	for {
		fmt.Print("Enter command (type 'help' for options): ")
		

		commands := io.GetInput()
		mainCommand := commands[0]
		args := commands[1:]
		switch mainCommand{
			case "help":
				fmt.Println("Available commands: help, stats, inventory, move, exit")
			case "stats":
				player.ShowDescription()
			case "inventory":
				player.ShowInventory()
			case "move":
				if len(args) < 1 {
					fmt.Println("Usage: move <direction>")
					continue
				}
				movementHandler(player, currentMap, args[0])
			case "map":
				currentMap.PrintMapWithPlayer(player)
				fmt.Printf("You are at position (%d, %d)\n", player.PlayerPositionY, player.PlayerPositionX)
				fmt.Printf("Current tile: %s\n", currentMap.Tiles[player.PlayerPositionY][player.PlayerPositionX].Name)
			case "exit":
				fmt.Println("Exiting game...")
				return
			default:
				fmt.Println("Unknown command. Type 'help' for options.")
		}
	}
}

func movementHandler(player *playerLogic.Player, currentMap *mapLogic.Map, direction string){
	switch direction{
		case "up", "w":
			_, err := currentMap.MovePlayer(player, -1, 0)
			if err != nil{
				fmt.Println("Move failed:", err)
			}
		case "down", "s":
			_, err := currentMap.MovePlayer(player, 1, 0)
			if err != nil{
				fmt.Println("Move failed:", err)
			}
		case "left", "a":
			_, err := currentMap.MovePlayer(player, 0, -1)
			if err != nil{
				fmt.Println("Move failed:", err)
			}
		case "right", "d":
			_, err := currentMap.MovePlayer(player, 0, 1)
			if err != nil{
				fmt.Println("Move failed:", err)
			}
		default:
			fmt.Println("Unknown direction. Use up, down, left, or right.")
	}
}