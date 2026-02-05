package mapLogic

import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/equipment"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
)

type Terrain int
const (
	TerrainGrass Terrain = iota
	TerrainWater
	TerrainMountain
	TerrainForest
)

type Encounter struct{
	NPCs []playerLogic.NPC `json:"npcs"`
}

type Tile struct{
	Name string `json:"name"`
	Terrain Terrain `json:"terrain"`
	Description string `json:"description"`
	Walkable bool `json:"walkable"`
	Encounter Encounter `json:"encounter,omitempty"`
	Entrance bool `json:"entrance,omitempty"`
	Exit bool `json:"exit,omitempty"`
	VisibleOnMap bool `json:"visible_on_map"`
	Equipment []equipment.Equipment `json:"equipment,omitempty"`
}

// GetDescription returns the description of the tile.
func (t *Tile) GetDescription() string{
	return t.Description
}

// GetName returns the name of the tile.
func (t *Tile) GetName() string{
	return t.Name
}

// Show provides a simple visual representation of the tile based on its terrain.
func (t *Tile) Display(){
	switch t.Terrain{
		case TerrainGrass:
			fmt.Print(".")
		case TerrainWater:
			fmt.Print("~")
		case TerrainMountain:
			fmt.Print("^")
		case TerrainForest:
			fmt.Print("*")
		default:
			fmt.Print("?")
		}
}

// Details prints detailed information about the tile.
func (t *Tile) Details(){
	fmt.Printf("Tile: %s\nTerrain: %d\nDescription: %s\nWalkable: %t\nEntrance: %t\nExit: %t\nVisible: %t", t.Name, t.Terrain, t.Description, t.Walkable, t.Entrance, t.Exit, t.VisibleOnMap)
	
	if len(t.Equipment) > 0{
		fmt.Println("\nItems on Tile:")
		for _, item := range t.Equipment{
			fmt.Printf("- %s\n", item.GetName())
		}
	}
	
	if len(t.Encounter.NPCs) > 0{
		fmt.Println("\nNPCs on Tile:")
		for _, npc := range t.Encounter.NPCs{
			fmt.Printf("- %s\n", npc.Name)
		}
	}
}