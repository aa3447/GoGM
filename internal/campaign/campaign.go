package campaign

import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/gameState"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
)

// Campaign represents a game campaign managed by a Game Master (GM).
type Campaign struct {
	Name string `json:"name"`
	Description string `json:"description"`
	GM playerLogic.GM `json:"gm"`
	Players map[string]*playerLogic.Player `json:"players"`
	NPCs map[string]*playerLogic.NPC `json:"npcs"`
	Monsters map[string]*playerLogic.NPC `json:"monsters"`
	Maps map[string]*mapLogic.Map `json:"maps"`
	CurrentMap *mapLogic.Map `json:"current_map,omitempty"`
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

	gameMap.GameState = &gameState.Gamestate{
		Players: c.Players,
	}

	c.Maps[gameMap.Name] = &gameMap

	return &gameMap, nil
}

func (c *Campaign) NewGamestateWithExistingMap(existingMap *mapLogic.Map) (*mapLogic.Map, error){
	existingMap.GameState = &gameState.Gamestate{
		Players: c.Players,
	}

	c.Maps[existingMap.Name] = existingMap

	return existingMap, nil
}

func (c *Campaign) AddMap(newMap *mapLogic.Map){
	c.Maps[newMap.Name] = newMap
}

func (c *Campaign) RemoveMap(mapName string){
	delete(c.Maps, mapName)
}

func (c *Campaign) SetCurrentMap(mapName string) error{
	selectedMap, exists := c.Maps[mapName]
	if !exists{
		return  fmt.Errorf("map %s does not exist in campaign", mapName)
	}
	c.CurrentMap = selectedMap
	return nil
}

func (c *Campaign) AddPlayer(newPlayer *playerLogic.Player){
	c.Players[newPlayer.Name] = newPlayer
}

func (c *Campaign) RemovePlayer(playerName string){
	delete(c.Players, playerName)
}

func (c *Campaign) GetPlayer(playerName string) (*playerLogic.Player, bool){
	player, exists := c.Players[playerName]
	return player, exists
}

func (c *Campaign) ListPlayers() []string{
	playerNames := []string{}
	for name := range c.Players{
		playerNames = append(playerNames, name)
	}
	return playerNames
}

func (c *Campaign) AddNPC(newNPC *playerLogic.NPC){
	c.NPCs[newNPC.Name] = newNPC
}

func (c *Campaign) RemoveNPC(npcName string){
	delete(c.NPCs, npcName)
}

func (c *Campaign) GetNPC(npcName string) (*playerLogic.NPC, bool){
	npc, exists := c.NPCs[npcName]
	return npc, exists
}

func (c *Campaign) ListNPCs() []string{
	npcNames := []string{}
	for name := range c.NPCs{
		npcNames = append(npcNames, name)
	}
	return npcNames
}

func (c *Campaign) SetGM(newGM playerLogic.GM){
	c.GM = newGM
}

func (c *Campaign) isJsonSafe(){}
