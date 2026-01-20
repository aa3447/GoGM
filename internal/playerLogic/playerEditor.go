package playerLogic

import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
)

func PlayerEditor(){
	commands := io.GetInput()
	basicPlayerVariables := []string{"name", "description", "background"}

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
				variable := args[0]
				value := args[1]

				/*player, exists := GetPlayerByID(playerID)
				if !exists {
					io.Print("Player not found")
					continue
				}
				currentPlayer = player*/
				
				// Check if variable is a basic player variable
				isBasicVariable := false
				for _, v := range basicPlayerVariables {
					if variable == v {
						isBasicVariable = true
						break
					}
				}
				if isBasicVariable {
					currentPlayer.editBasicPlayerVariables(variable, value)
				} else {
					fmt.Println("Unknown variable")
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