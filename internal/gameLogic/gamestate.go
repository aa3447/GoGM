package gameLogic

import (
	"fmt"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
)

type Gamestate struct{
	CurrentMap [][]mapLogic.Tile
	PlayerPositionX int
	PlayerPositionY int
}

func NewGamestateWithRandomMap(mapSizeY int, mapSizeX int ,encounterProbability float64 ,terrainProbability []float64) (*Gamestate, error){
	gameMap ,entranceLocation ,err := mapLogic.GenRandomMap(mapSizeY, mapSizeX, encounterProbability, terrainProbability)
	if err != nil{
		return nil, err
	}
	return &Gamestate{
		CurrentMap: gameMap,
		PlayerPositionX: entranceLocation[1],
		PlayerPositionY: entranceLocation[0],
	}, nil
}

func (gs *Gamestate) MovePlayer(deltaY int, deltaX int) error{
	newY := gs.PlayerPositionY + deltaY
	newX := gs.PlayerPositionX + deltaX
	if newY < 0 || newY >= len(gs.CurrentMap) || newX < 0 || newX >= len(gs.CurrentMap[0]){
		fmt.Println("Cannot move out of bounds")
		return nil // Out of bounds, ignore move
	}
	if !gs.CurrentMap[newY][newX].Walkable{
		fmt.Println("Tile not walkable")
		return nil // Tile not walkable, ignore move
	}
	gs.PlayerPositionY = newY
	gs.PlayerPositionX = newX
	gs.CurrentMap[newY][newX].VisibleOnMap = true
	return nil
}

func (gs *Gamestate) GetPlayerPosition() (int, int){
	return gs.PlayerPositionY, gs.PlayerPositionX
}