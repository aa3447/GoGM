package mapLogic



type Terrain int
const (
	TerrainGrass Terrain = iota
	TerrainWater
	TerrainMountain
	TerrainForest
)

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

func (t *Tile) GetDescription() string{
	return t.Description
}

func (t *Tile) GetName() string{
	return t.Name
}

