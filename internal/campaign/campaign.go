package campaign

import (
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/gameLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
)

type Campaign struct {
	Name string
	Description string
	GM playerLogic.GM
	Players map[string]*playerLogic.Player
	Maps map[string]*mapLogic.Map
}

func NewCampaign(name string, description string, gm playerLogic.GM) *Campaign{
	return &Campaign{
		Name: name,
		Description: description,
		GM: gm,
		Players: map[string]*playerLogic.Player{},
		Maps: map[string]*mapLogic.Map{},
	}
}

func (c *Campaign) NewGamestateWithRandomMap(mapSizeY int, mapSizeX int ,encounterProbability float64 ,terrainProbability []float64, mapName string) (*mapLogic.Map, error){
	gameMap , err := mapLogic.GenRandomMap(mapSizeY, mapSizeX, encounterProbability, terrainProbability,mapName)
	if err != nil{
		return nil ,err
	}

	gameMap.GameState = &gameLogic.Gamestate{
		Players: c.Players,
	}

	c.Maps[gameMap.Name] = &gameMap

	return &gameMap, nil
}

func (c *Campaign) NewGamestateWithExistingMap(existingMap *mapLogic.Map) (*mapLogic.Map, error){
	existingMap.GameState = &gameLogic.Gamestate{
		Players: c.Players,
	}

	c.Maps[existingMap.Name] = existingMap

	return existingMap, nil
}

func (c *Campaign) AddMap(newMap *mapLogic.Map){
	c.Maps[newMap.Name] = newMap
}

func (c *Campaign) AddPlayer(newPlayer *playerLogic.Player){
	c.Players[newPlayer.Name] = newPlayer
}

func (c *Campaign) RemovePlayer(playerName string){
	delete(c.Players, playerName)
}

func (c *Campaign) SetGM(newGM playerLogic.GM){
	c.GM = newGM
}
