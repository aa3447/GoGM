package gameLogic

import (
	"fmt"
	"log"
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

func (gs *Gamestate) MovePlayer(player *player.Player ,deltaY int, deltaX int) (mapLogic.PlayerMove, error){
	newY := player.PlayerPositionY + deltaY
	newX := player.PlayerPositionX + deltaX
	CurrentMapTiles := gs.CurrentMap.Tiles
	if newY < 0 || newY >= len(CurrentMapTiles) || newX < 0 || newX >= len(CurrentMapTiles[0]){
		return mapLogic.PlayerMove{}, fmt.Errorf("cannot move out of bounds") // Out of bounds, ignore move
	}
	if !CurrentMapTiles[newY][newX].Walkable{
		return mapLogic.PlayerMove{},fmt.Errorf("tile not walkable") // Tile not walkable, ignore move
	}
	player.SetLocation(newY, newX)
	CurrentMapTiles[newY][newX].VisibleOnMap = true

	playerMove := mapLogic.PlayerMove{
		PlayerName: player.Name,
		From: []int{player.PlayerPositionY - deltaY, player.PlayerPositionX - deltaX},
		To: []int{newY, newX},
	}
	log.Printf("Player Move: %+v\n", playerMove)

	return playerMove,nil
}


func (gs *Gamestate) AddMap(newMap *mapLogic.Map){
	gs.Maps = append(gs.Maps, newMap)
}
