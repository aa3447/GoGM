package mapLogic

import (
	"errors"
	"fmt"
	player "home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"math"
	"math/rand"
	"slices"
)

type NPC struct{
	Name string
	Description string
	Hostile bool
	Health int
	Attack int
	Defense int
}

type Encounter struct{
	NPCs []NPC `json:"npcs"`
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
}

type Map struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Tiles [][]Tile `json:"tiles"`
}

type Terrain int
const (
	TerrainGrass Terrain = iota
	TerrainWater
	TerrainMountain
	TerrainForest
)

// Generate a random map with given sizeX and terrain probabilities
//Terrain probabilities go in the order of Grass, Water, Mountain, Forest and should sum to 1.0
func GenRandomMap(sizeY int, sizeX int, encounterProbability float64 ,terrainProbability []float64) (Map, []int, error){
	if sizeY <= 0 || sizeX <= 0{
		return Map{}, nil, errors.New("invalid map sizeX")
	}
	if len(terrainProbability) != 4{
		return Map{}, nil, errors.New("Terrain probability must have exactly 4 elements")
	}
	if math.Abs(terrainProbability[0] + terrainProbability[1] + terrainProbability[2] + terrainProbability[3] - 1.0) > 0.0001{
		return Map{}, nil, errors.New("Terrain probabilities must sum to 1.0")
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

	entranceLocation := generateEntranceAndExit(tileMap)

	newMap := Map{
		Name: "Random Map",
		Description: "A randomly generated map.",
		Tiles: tileMap,
	}

	return newMap, entranceLocation, nil
}

func generateEntranceAndExit(tileMap [][]Tile) []int{
	sizeY := len(tileMap)
	sizeX := len(tileMap[0])
	entranceLocation := []int{}

	topOrSideEntrance := rand.Float64()
	if topOrSideEntrance < 0.5{
		entranceY := rand.Intn(sizeY)
		exitY := rand.Intn(sizeY)
		leftOrRightEntrance := int(math.Round(rand.Float64())) * (sizeX - 1)
		tileMap[entranceY][leftOrRightEntrance].Entrance = true
		tileMap[entranceY][leftOrRightEntrance].Walkable = true
		tileMap[entranceY][leftOrRightEntrance].VisibleOnMap = true
		entranceLocation = append(entranceLocation, entranceY, leftOrRightEntrance)
		if leftOrRightEntrance == 0{{
			tileMap[exitY][sizeX-1].Exit = true
			tileMap[exitY][sizeX-1].Walkable = true
		}
		} else {
			tileMap[exitY][0].Exit = true
			tileMap[exitY][0].Walkable = true
		}
	} else {
		entranceX := rand.Intn(sizeX)
		exitX := rand.Intn(sizeX)
		topOrBottomEntrance := int(math.Round(rand.Float64())) * (sizeY - 1)
		tileMap[topOrBottomEntrance][entranceX].Entrance = true
		tileMap[topOrBottomEntrance][entranceX].Walkable = true
		tileMap[topOrBottomEntrance][entranceX].VisibleOnMap = true
		entranceLocation = append(entranceLocation, topOrBottomEntrance, entranceX)
		if topOrBottomEntrance == 0{
			tileMap[sizeY-1][exitX].Exit = true
			tileMap[sizeY-1][exitX].Walkable = true
		} else {
			tileMap[0][exitX].Exit = true
			tileMap[0][exitX].Walkable = true
		}
	}

	return  entranceLocation
}

func PrintMap(tileMap [][]Tile){
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

func PrintMapDebug(tileMap [][]Tile){
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

func PrintMapWithPlayer(tileMap [][]Tile, player *player.Player){
	playerY, playerX := player.GetPlayerPosition()
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
