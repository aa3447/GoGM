package editors

import (
	"fmt"
	"strconv"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"
)

// MapEditor provides an interface to create and edit maps.
func MapEditor(){
	var currentMap *mapLogic.Map

	for {
		commands := io.GetInput()
		command := commands[0]
		args := commands[1:]

		switch command{
			case "create", "c":
				switch args[0]{
					case "random":
						//currentMap = mapLogic.GenerateRandomMap(args[1], args[2])
					case "empty":
						y, err := strconv.Atoi(args[1])
						if err != nil {
							fmt.Println("Error parsing map height:", err)
							continue
						}
						x, err := strconv.Atoi(args[2])
						if err != nil {
							fmt.Println("Error parsing map width:", err)
							continue
						}
						currentMap, err = mapLogic.GenEmptyMap(y, x, args[3])
						if err != nil {
							fmt.Println("Error creating empty map:", err)
						} else {
							fmt.Println("Empty map created successfully.")
						}
					default:
						fmt.Println("Unknown map type. Available types: random, empty")
				}
			case "edit", "e":
				switch args[0]{
					case "name":
						if currentMap == nil{
							fmt.Println("No map loaded. Please create a map first.")
							continue
						}
						currentMap.Name = args[1]
					case "description":
						if currentMap == nil{
							fmt.Println("No map loaded. Please create a map first.")
							continue
						}
						currentMap.Description = args[1]
					default:
						fmt.Println("Unknown edit command. Available commands: name, description")
				}
			case "save", "s":
				var saveName string
				if len(args) >= 1{
					saveName = args[0]
				} else {
					saveName = currentMap.Name
				}
				if err := serialization.SaveToFile(*currentMap, "gm", "map", saveName);  err != nil {
					fmt.Println("Error saving map:", err)
				} else {
					fmt.Println("Map saved successfully.")
				}
			case "load", "l":
				if len(args) < 1{
					fmt.Println("Usage: load  <map_name>")
					continue
				}
				loadedMap, err := serialization.LoadFromJSONFile("gm", "map", args[0], mapLogic.Map{})
				if err != nil{
					fmt.Println("Error loading map:", err)
				} else {
					currentMap = loadedMap
					fmt.Println("Map loaded successfully.")
				}
			case "show", "view", "v":
				if currentMap == nil{
					fmt.Println("No map loaded. Please create or load a map first.")
					continue
				}
				fmt.Printf("Map Name: %s\nDescription: %s\nSize: %dx%d\n", currentMap.Name, currentMap.Description, len(currentMap.Tiles), len(currentMap.Tiles[0]))
				currentMap.PrintMapDebug()
			default:
				fmt.Println("Unknown command:", command)
		}
	}
}

// tileEditor provides an interface to edit tiles within a map.
func tileEditor(currentMap *mapLogic.Map) *mapLogic.Map{
	var err error
	currentTile := &currentMap.Tiles[0][0]
	
	for {
		fmt.Printf("Editing Tile: %s\n", currentTile.Name)

		commands := io.GetInput()
		command := commands[0]
		args := commands[1:]

		switch command{
			case "set", "edit":
				currentTile, err = tileOptionsEditor(currentTile, args[0], args[1])
				if err != nil {
					fmt.Println("Error setting tile option:", err)
				}
			case "show":
				// reserved for tile showing
			case "select":
				y, err := strconv.Atoi(args[0])
				if err != nil {
					fmt.Println("Error parsing tile Y coordinate:", err)
					continue
				}
				x, err := strconv.Atoi(args[1])
				if err != nil {
					fmt.Println("Error parsing tile X coordinate:", err)
					continue
				}
				if y < 0 || y >= len(currentMap.Tiles) || x < 0 || x >= len(currentMap.Tiles[0]) {
					fmt.Println("Tile coordinates are out of bounds.")
					continue
				}
				currentTile = &currentMap.Tiles[y][x]
			case "exit":
				return currentMap
			default:
				fmt.Println("Unknown command:", command)
		}
	}
}

// tileOptionsEditor sets options for a specific tile.
func tileOptionsEditor(currentTile *mapLogic.Tile, what, value string) (*mapLogic.Tile, error){
	var err error
	fmt.Printf("Editing Tile Options for: %s\n", currentTile.Name)

	switch what{
		case "name":
			currentTile.Name = value
		case "terrain":
			//
		case "description":
			currentTile.Description = value
		case "walkable":
			currentTile, err = tileBooleanSelector(currentTile, "walkable", value)
		case "encounter":
			//
		case "entrance":
			currentTile, err = tileBooleanSelector(currentTile, "entrance", value)
		case "exit":
			currentTile, err = tileBooleanSelector(currentTile, "exit", value)
		case "visible":
			currentTile, err = tileBooleanSelector(currentTile, "visible", value)
		case "items":
			//	
		default:
			fmt.Println("Unknown set command:", what)
	}

	if err != nil{
		return currentTile, fmt.Errorf("failed to set tile option %s: %v", what, err)
	}
	return currentTile, nil
}

// tileBooleanSelector helps set boolean properties of a tile based on user input.
func tileBooleanSelector(currentTile *mapLogic.Tile, command, value string) (*mapLogic.Tile, error){
	
	var boolValue bool
	switch value {
		case "yes", "y", "true":
			boolValue = true
		case "no", "n", "false":
			boolValue = false
		default:
			return currentTile, fmt.Errorf("invalid boolean value: %s", value)
		}

	switch command{
		case "walkable":
			currentTile.Walkable = boolValue
		case "entrance":
			currentTile.Entrance = boolValue
		case "exit":
			currentTile.Exit = boolValue
		case "visible":
			currentTile.VisibleOnMap = boolValue
		default:
			return currentTile, fmt.Errorf("unknown command: %s", command)
	}
	return currentTile, nil
}