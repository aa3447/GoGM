package gameLogic

import (
	"fmt"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
)

type Gamestate struct{
	CurrentMap mapLogic.Map
}

func NewGamestateWithRandomMap(mapSizeY int, mapSizeX int ,encounterProbability float64 ,terrainProbability []float64) (*Gamestate, []int, error){
	gameMap ,entranceLocation ,err := mapLogic.GenRandomMap(mapSizeY, mapSizeX, encounterProbability, terrainProbability)
	if err != nil{
		return nil, []int{},err
	}
	return &Gamestate{
		CurrentMap: gameMap,
	}, entranceLocation, nil
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

