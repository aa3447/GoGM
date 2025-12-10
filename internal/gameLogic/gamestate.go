package gameLogic

import (
	"fmt"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
)

type Gamestate struct{
	CurrentMap *mapLogic.Map
	Maps []*mapLogic.Map
	Players []*player.Player
}

func NewGamestateWithRandomMap(mapSizeY int, mapSizeX int ,encounterProbability float64 ,terrainProbability []float64, mapName string) (*Gamestate, error){
	gameMap , err := mapLogic.GenRandomMap(mapSizeY, mapSizeX, encounterProbability, terrainProbability,mapName)
	if err != nil{
		return nil ,err
	}
	return &Gamestate{
		CurrentMap: &gameMap,
		Maps: []*mapLogic.Map{&gameMap},
	},  nil
}

func NewGamestateWithExistingMap(existingMap *mapLogic.Map) *Gamestate{
	return &Gamestate{
		CurrentMap: existingMap,
		Maps: []*mapLogic.Map{existingMap},
	}
}

func (gs *Gamestate) MovePlayer(player *player.Player ,deltaY int, deltaX int) error{
	newY := player.PlayerPositionY + deltaY
	newX := player.PlayerPositionX + deltaX
	CurrentMapTiles := gs.CurrentMap.Tiles
	if newY < 0 || newY >= len(CurrentMapTiles) || newX < 0 || newX >= len(CurrentMapTiles[0]){
		fmt.Println("Cannot move out of bounds")
		return nil // Out of bounds, ignore move
	}
	if !CurrentMapTiles[newY][newX].Walkable{
		fmt.Println("Tile not walkable")
		return nil // Tile not walkable, ignore move
	}
	player.SetLocation(newY, newX)
	CurrentMapTiles[newY][newX].VisibleOnMap = true
	return nil
}


func (gs *Gamestate) AddMap(newMap *mapLogic.Map){
	gs.Maps = append(gs.Maps, newMap)
}
