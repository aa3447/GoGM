package mapLogic

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"log"

	player "home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/gameState"
)

type Map struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Tiles [][]Tile `json:"tiles"`
	EntranceLocation []int `json:"entrance_location,omitempty"`
	ExitLocation []int `json:"exit_location,omitempty"`
	FileLocation string `json:"file_location,omitempty"`
	GameState *gameState.Gamestate `json:"-"`

}

type PlayerMove struct{
	PlayerName string `json:"player_name"`
	From []int `json:"from"`
	To []int `json:"to"`
}

// Generate a random map with given sizeX and terrain probabilities
//Terrain probabilities go in the order of Grass, Water, Mountain, Forest and should sum to 1.0
func GenRandomMap(sizeY int, sizeX int, encounterProbability float64 ,terrainProbability []float64, name string) (Map,error){
	if sizeY <= 0 || sizeX <= 0{
		return Map{}, errors.New("invalid map sizeX")
	}
	if len(terrainProbability) != 4{
		return Map{}, errors.New("Terrain probability must have exactly 4 elements")
	}
	if math.Abs(terrainProbability[0] + terrainProbability[1] + terrainProbability[2] + terrainProbability[3] - 1.0) > 0.0001{
		return Map{}, errors.New("Terrain probabilities must sum to 1.0")
	}

	slices.Sort(terrainProbability)

	tileMap := make([][]Tile, sizeY)
	for rowIndex := range sizeY{
		for columnIndex := range sizeX{
			tile := Tile{
				Name: fmt.Sprintf("%d,%d", rowIndex, columnIndex),
				VisibleOnMap: false,
			}

			randVal := rand.Float64()
			switch {
				case randVal < terrainProbability[0]:
					tile.Terrain = TerrainGrass
					tile.Description = "A grassy field."
					tile.Walkable = true
				case randVal < terrainProbability[1]:
					tile.Terrain = TerrainWater
					tile.Description = "A body of water."
					tile.Walkable = false
				case randVal < terrainProbability[2]:
					tile.Terrain = TerrainMountain
					tile.Description = "A towering mountain."
					tile.Walkable = false
				default:
					tile.Terrain = TerrainForest
					tile.Description = "A dense forest."
					tile.Walkable = true
			}

			encounterRoll := rand.Float64()
			if encounterRoll < encounterProbability{
				//tile.Encounter = generateRandomEncounter()
				tile.Encounter = Encounter{}
			}

			tileMap[rowIndex] = append(tileMap[rowIndex], tile)
		}
	}

	entranceAndExitLocation := generateEntranceAndExit(tileMap)

	if name == ""{
		name = "Random Map"
	}
	newMap := Map{
		Name: name,
		Description: "A randomly generated map.",
		Tiles: tileMap,
		EntranceLocation: entranceAndExitLocation[0:2],
		ExitLocation: entranceAndExitLocation[2:4],
	}

	return newMap, nil
}

// generateEntranceAndExit randomly places entrance and exit points on the map edges
func generateEntranceAndExit(tileMap [][]Tile) []int{
	sizeY := len(tileMap)
	sizeX := len(tileMap[0])
	entranceAndExitLocation := []int{}

	topOrSideEntrance := rand.Float64()
	if topOrSideEntrance < 0.5{
		entranceY := rand.Intn(sizeY)
		exitY := rand.Intn(sizeY)
		leftOrRightEntrance := int(math.Round(rand.Float64())) * (sizeX - 1)
		tileMap[entranceY][leftOrRightEntrance].Entrance = true
		tileMap[entranceY][leftOrRightEntrance].Walkable = true
		tileMap[entranceY][leftOrRightEntrance].VisibleOnMap = true
		entranceAndExitLocation = append(entranceAndExitLocation, entranceY, leftOrRightEntrance)
		if leftOrRightEntrance == 0{{
			tileMap[exitY][sizeX-1].Exit = true
			tileMap[exitY][sizeX-1].Walkable = true
			entranceAndExitLocation = append(entranceAndExitLocation, exitY, sizeX-1)
		}
		} else {
			tileMap[exitY][0].Exit = true
			tileMap[exitY][0].Walkable = true
			entranceAndExitLocation = append(entranceAndExitLocation, exitY, 0)
		}
	} else {
		entranceX := rand.Intn(sizeX)
		exitX := rand.Intn(sizeX)
		topOrBottomEntrance := int(math.Round(rand.Float64())) * (sizeY - 1)
		tileMap[topOrBottomEntrance][entranceX].Entrance = true
		tileMap[topOrBottomEntrance][entranceX].Walkable = true
		tileMap[topOrBottomEntrance][entranceX].VisibleOnMap = true
		entranceAndExitLocation = append(entranceAndExitLocation, topOrBottomEntrance, entranceX)
		if topOrBottomEntrance == 0{
			tileMap[sizeY-1][exitX].Exit = true
			tileMap[sizeY-1][exitX].Walkable = true
			entranceAndExitLocation = append(entranceAndExitLocation, sizeY-1, exitX)
		} else {
			tileMap[0][exitX].Exit = true
			tileMap[0][exitX].Walkable = true
			entranceAndExitLocation = append(entranceAndExitLocation, 0, exitX)
		}
	}

	return  entranceAndExitLocation
}

// GenEmptyMap generates an empty map with all grassy tiles
func GenEmptyMap(sizeY int, sizeX int, name string) (*Map, error){
	if sizeY <= 0 || sizeX <= 0{
		return nil, errors.New("invalid map sizeX")
	}

	tileMap := make([][]Tile, sizeY)
	for rowIndex := range sizeY{
		for columnIndex := range sizeX{
			tile := Tile{
				Name: fmt.Sprintf("%d,%d", rowIndex, columnIndex),
				Terrain: TerrainGrass,
				Description: "An empty grassy field.",
				Walkable: true,
				VisibleOnMap: false,
			}
			tileMap[rowIndex] = append(tileMap[rowIndex], tile)
		}
	}

	if name == ""{
		name = "Empty Map"
	}
	m := Map{
		Name: name,
		Description: "An empty map.",
		Tiles: tileMap,
		EntranceLocation: []int{},
		ExitLocation: []int{},
	}

	return &m, nil
}

// PrintMap prints the map to the console, showing only visible tiles.
func (m *Map) PrintMap(){
	tileMap := m.Tiles
	for _,row := range tileMap{
		for _,tile := range row{
			if tile.VisibleOnMap{
				switch tile.Terrain{
					case TerrainGrass:
						fmt.Print(".")
					case TerrainWater:
						fmt.Print("~")
					case TerrainMountain:
						fmt.Print("^")
					case TerrainForest:
						fmt.Print("*")
					}
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

// PrintMapDebug prints the entire map to the console, showing all tiles regardless of visibility.
func (m *Map) PrintMapDebug(){
	tileMap := m.Tiles
	for _,row := range tileMap{
		for _,tile := range row{
			switch tile.Terrain{
				case TerrainGrass:
					fmt.Print(".")
				case TerrainWater:
					fmt.Print("~")
				case TerrainMountain:
					fmt.Print("^")
				case TerrainForest:
					fmt.Print("*")
			}
		}
		fmt.Println()
	}
}

// PrintMapDebugWithPlayers prints the entire map to the console, showing all tiles and player positions.
func (m *Map) PrintMapDebugWithPlayers(players []*player.Player){
	tileMap := m.Tiles
	for rowIndex, row := range tileMap{
		for columnIndex, tile := range row{
			playerHere := false
			for _, p := range players{
				playerY, playerX := p.GetPlayerPosition()
				if playerY == rowIndex && playerX == columnIndex{
					playerHere = true
					break
				}
			}
			if playerHere{
				fmt.Print("P")
			} else {
				switch tile.Terrain{
					case TerrainGrass:
						fmt.Print(".")
					case TerrainWater:
						fmt.Print("~")
					case TerrainMountain:
						fmt.Print("^")
					case TerrainForest:
						fmt.Print("*")
				}
			}
		}
		fmt.Println()
	}
}

// PrintMapWithPlayer prints the map to the console, showing only visible tiles and the specified player's position.
func (m *Map) PrintMapWithPlayer(player *player.Player){
	playerY, playerX := player.GetPlayerPosition()
	tileMap := m.Tiles
	for rowIndex, row := range tileMap{
		for columnIndex, tile := range row{
			if rowIndex == playerY && columnIndex == playerX{
				fmt.Print("P")
			} else if tile.VisibleOnMap{
				switch tile.Terrain{
					case TerrainGrass:
						fmt.Print(".")
					case TerrainWater:
						fmt.Print("~")
					case TerrainMountain:
						fmt.Print("^")
					case TerrainForest:
						fmt.Print("*")
					}
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

// PrintMapWithPlayers prints the map to the console, showing only visible tiles and the players position.
func (m *Map) PrintMapWithPlayers(players []*player.Player){
	tileMap := m.Tiles
	for rowIndex, row := range tileMap{
		for columnIndex, tile := range row{
			playerHere := false
			for _, p := range players{
				playerY, playerX := p.GetPlayerPosition()
				if playerY == rowIndex && playerX == columnIndex{
					playerHere = true
					break
				}
			}
			if playerHere{
				fmt.Print("P")
			} else if tile.VisibleOnMap{
				switch tile.Terrain{
					case TerrainGrass:
						fmt.Print(".")
					case TerrainWater:
						fmt.Print("~")
					case TerrainMountain:
						fmt.Print("^")
					case TerrainForest:
						fmt.Print("*")
				}
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

// MovePlayer attempts to move the player by deltaY and deltaX. Returns a PlayerMove struct and error if move is invalid.
func (m *Map) MovePlayer(player *player.Player ,deltaY int, deltaX int) (PlayerMove, error){
	newY := player.PlayerPositionY + deltaY
	newX := player.PlayerPositionX + deltaX
	MapTiles := m.Tiles
	if newY < 0 || newY >= len(MapTiles) || newX < 0 || newX >= len(MapTiles[0]){
		return PlayerMove{}, fmt.Errorf("cannot move out of bounds") // Out of bounds, ignore move
	}
	if !MapTiles[newY][newX].Walkable{
		return PlayerMove{},fmt.Errorf("tile not walkable") // Tile not walkable, ignore move
	}
	player.SetLocation(newY, newX)
	MapTiles[newY][newX].VisibleOnMap = true

	playerMove := PlayerMove{
		PlayerName: player.Name,
		From: []int{player.PlayerPositionY - deltaY, player.PlayerPositionX - deltaX},
		To: []int{newY, newX},
	}
	log.Printf("Player Move: %+v\n", playerMove)

	return playerMove,nil
}

func (m *Map) GetTileAt(y int, x int) (Tile, error){
	if y < 0 || y >= len(m.Tiles) || x < 0 || x >= len(m.Tiles[0]){
		return Tile{}, fmt.Errorf("coordinates out of bounds")
	}
	return m.Tiles[y][x], nil
}

func (m *Map) SetTileAt(y int, x int, tile Tile) error{
	if y < 0 || y >= len(m.Tiles) || x < 0 || x >= len(m.Tiles[0]){
		return fmt.Errorf("coordinates out of bounds")
	}
	m.Tiles[y][x] = tile
	return nil
}


