package playerEditor

import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"
)

func PlayerEditor(player ...*playerLogic.Player) *playerLogic.Player{
	var currentPlayer *playerLogic.Player
	if len(player) > 0{
		currentPlayer = player[0]
	}
	
	for {
		commands := io.GetInput()
		command := commands[0]
		args := commands[1:]
				
		switch command{
			case "create":
				currentPlayer = playerLogic.CreatePlayer()
			case "edit":
				if len(args) < 2 {
					fmt.Println("Usage: edit  <variable> <value>")
					continue
				}
				if currentPlayer == nil{
					fmt.Println("No player loaded. Please create a player first.")
					continue
				}
				whatToEdit := args[0]
				value := args[1]

				switch whatToEdit{
					case "name", "description", "background":
						currentPlayer.EditBasicPlayerVariables(whatToEdit, value)
					case "class":
						if err := currentPlayer.ChangeClass(value); err != nil {
							fmt.Println("Error changing class:", err)
						}
					case "level":
						var level int
						fmt.Sscanf(value, "%d", &level)
						
						if err := currentPlayer.ChangeLevel(level); err != nil {
							fmt.Println("Error setting level:", err)
						}
					case "level_track":
						levelTrack, err := playerLogic.GetLevelTrack(value)
						if err != nil {
							fmt.Println("Error changing level track:", err)
							continue
						}
						if err := currentPlayer.ChangeLevelTrack(levelTrack); err != nil {
							fmt.Println("Error changing level track:", err)
						}
					case "exp":
						var exp int
						fmt.Sscanf(value, "%d", &exp)
						if err := currentPlayer.ChangeExperience(exp); err != nil {
							fmt.Println("Error changing experience:", err)
						}
					case "attribute":
						if len(args) < 3 {
							fmt.Println("Usage: edit attribute <attribute_name> <value>")
							continue
						}
						attributeName := args[1]
						if attributeName == "all"{
							if len(args) < 8 {
								fmt.Println("Usage: edit attribute all <strength> <dexterity> <intelligence> <constitution> <charisma> <wisdom>")
								continue
							}
							var attributeValues [6]int
							for index := range currentPlayer.Attributes.ToSlice(){
								var attributeValue int
								fmt.Sscanf(args[index+2], "%d", &attributeValue)
								attributeValues[index] = attributeValue
							}
							if err := currentPlayer.ChangeAllAttributes(attributeValues); err != nil {
								fmt.Println("Error changing attributes:", err)
							}
						} else {
							var attributeValue int
							fmt.Sscanf(args[2], "%d", &attributeValue)
							if err := currentPlayer.ChangeAttribute(attributeName, attributeValue); err != nil {
								fmt.Println("Error changing attribute:", err)
							}
						}
					case "inventory":
						// Reserve for inventory manager
					case "hp":
						var how string
						fmt.Sscanf(value, "%s", &how)
						switch how {
							case "set":
								if len(args) < 3 {
									fmt.Println("Usage: edit hp set <value>")
									continue
								}
								var hpValue int
								fmt.Sscanf(args[2], "%d", &hpValue)
								if err := currentPlayer.SetHP(hpValue); err != nil {
									fmt.Println("Error setting HP:", err)
								}
							case "roll":
								currentPlayer.RerollAllHitDice()
							}
					default:
						fmt.Println("Unknown variable to edit:", whatToEdit)
				}
			case "save":
				if currentPlayer == nil{
					fmt.Println("No player loaded. Please create a player first.")
					continue
				}
				
				var saveName string
				if len(args) >= 1{
					saveName = args[0]
				} else {
					saveName = currentPlayer.Name
				}
				if err := serialization.SaveToFile(*currentPlayer, "player", "player", saveName);  err != nil {
					fmt.Println("Error saving player:", err)
				} else {
					fmt.Println("Player saved successfully.")
				}
			case "load":
				if len(args) < 1{
					fmt.Println("Usage: load  <player_name>")
					continue
				}
				loadedPlayer, err := serialization.LoadFromJSONFile("player", "player", args[0], playerLogic.Player{})
				if err != nil {
					fmt.Println("Error loading player:", err)
					continue
				}
				currentPlayer = loadedPlayer
			case "exit":
				fmt.Println("Exiting Player Editor.")
				return currentPlayer
			default:
				fmt.Println("Unknown command:", command)
		}
	}
}
