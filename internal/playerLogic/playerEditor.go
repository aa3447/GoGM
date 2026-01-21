package playerLogic

import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
)

func PlayerEditor(){
	commands := io.GetInput()

	for {
		command := commands[0]
		args := commands[1:]
		
		var currentPlayer *Player
		switch command{
			case "create":
				currentPlayer = CreatePlayer()
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
						currentPlayer.editBasicPlayerVariables(whatToEdit, value)
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
						levelTrack, err := GetLevelTrack(value)
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
					default:
						fmt.Println("Unknown variable to edit:", whatToEdit)
				}
		}
	}
}

func (p *Player) editBasicPlayerVariables(variable, value string){
		switch variable{
			case "name":
				p.Name = value
			case "description":
				p.Description = value
			case "background":
				p.Background = value
		}
}