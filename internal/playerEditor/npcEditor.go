package playerEditor


import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"
)

func NpcEditor(){
	var currentNPC *playerLogic.NPC

	for {
		commands := io.GetInput()
		command := commands[0]
		args := commands[1:]
		
		switch command{
			case "create":
				if len(args) < 2{
					fmt.Println("Usage: create  <name> <description>")
					continue
				}
				currentNPC = playerLogic.NewNPC(args[0], args[1])
			case "edit":
				if len(args) < 2 {
					fmt.Println("Usage: edit  <variable> <value>")
					continue
				}
				if currentNPC == nil{
					fmt.Println("No NPC loaded. Please create an NPC first.")
					continue
				}
				whatToEdit := args[0]

				if whatToEdit == "dialogue"{
					// reserved for dialogue editor
				} else {
					currentNPC.Player = *PlayerEditor(&currentNPC.Player)
				}
			case "save":
				var saveName string
				if len(args) >= 1{
					saveName = args[0]
				} else {
					saveName = currentNPC.Name
				}
				if err := serialization.SaveToFile(*currentNPC, "npc", "npc", saveName);  err != nil {
					fmt.Println("Error saving NPC:", err)
				} else {
					fmt.Println("NPC saved successfully.")
				}
			default:
				fmt.Println("Unknown command:", command)
		}
	}
}