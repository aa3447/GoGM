package editors


import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"
)

// NpcEditor provides an interface to create and edit NPCs.
func NpcEditor(){
	var currentNPC *playerLogic.NPC

	for {
		commands := io.GetInput()
		command := commands[0]
		args := commands[1:]
		
		switch command{
			case "create", "c":
				if len(args) < 2{
					fmt.Println("Usage: create  <name> <description>")
					continue
				}
				currentNPC = playerLogic.NewNPC(args[0], args[1])
			case "edit", "e":
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
			case "save", "s":
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
			case "load", "l":
				if len(args) < 1{
					fmt.Println("Usage: load  <npc_name>")
					continue
				}
				loadedNPC, err := serialization.LoadFromJSONFile("gm", "npc", args[0], playerLogic.NPC{})
				if err != nil{
					fmt.Println("Error loading NPC:", err)
				} else {
					currentNPC = loadedNPC
					fmt.Println("NPC loaded successfully.")
				}
			case "view", "show", "v":
				// View NPC stats
			case "quit", "exit", "q":
				return
			default:
				fmt.Println("Unknown command:", command)
		}
	}
}