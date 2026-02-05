package equipment

import "fmt"

type ArmorType int
const (
	ArmorTypeLight ArmorType  = iota
	ArmorTypeMedium
	ArmorTypeHeavy
)


type Armor struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Defense int `json:"defense"`
	Type ArmorType `json:"type"`
	Weight int `json:"weight"`
	Rarity Rarity `json:"rarity"`
	IsCustom bool `json:"is_custom"`
}

// Predefined armor types
var mapOfArmors = map[string]Armor{
	"Leather Armor": {
		Name: "Leather Armor",
		Description: "Basic leather armor offering minimal protection.",
		Defense: 5,
		Type: ArmorTypeLight,
		Weight: 10,
		Rarity: Common,
		IsCustom: false,
	},
	"Chainmail": {
		Name: "Chainmail",
		Description: "Interlinked metal rings providing moderate protection.",
		Defense: 15,
		Type: ArmorTypeMedium,
		Weight: 30,
		Rarity: Uncommon,
		IsCustom: false,
	},
	"Plate Armor": {
		Name: "Plate Armor",
		Description: "Heavy metal plates offering high protection.",
		Defense: 25,
		Type: ArmorTypeHeavy,
		Weight: 50,
		Rarity: Rare,
		IsCustom: false,
	},
}

// GetArmorByName retrieves an Armor by its name from the predefined armor types.
func GetArmorByName(name string) (Armor, bool){
	armor, exists := mapOfArmors[name]
	return armor, exists
}

// GetAllArmorTypes returns a slice of all predefined armor types.
func GetAllArmorTypes() []Armor{
	armors := []Armor{}
	for _, armor := range mapOfArmors{
		armors = append(armors, armor)
	}
	return armors
}

// GetAllArmors returns a slice of all armors, including custom ones.
func GetAllArmors() []Armor {
	armors := []Armor{}
	for _, armor := range mapOfArmors{
		armors = append(armors, armor)
	}
	return armors
}

// AddCustomArmor adds a new custom armor to the armor map.
func AddCustomArmor(name string, description string, defense int, t ArmorType, weight int, rarity Rarity){
	if !verifyArmorType(t){
		fmt.Println("Invalid armor type. Armor not added.")
		return
	}
	if !verifyRarity(rarity){
		fmt.Println("Invalid rarity. Armor not added.")
		return
	}

	mapOfArmors[name] = Armor{
		Name: name,
		Description: description,
		Defense: defense,
		Type: t,
		Weight: weight,
		Rarity: rarity,
		IsCustom: true,
	}

	fmt.Println("Custom armor added:", name)
}

// RemoveCustomArmor removes a custom armor from the armor map.
func removeCustomArmor(name string){
	armor, exists := mapOfArmors[name]
	if exists && armor.IsCustom{
		delete(mapOfArmors, name)
		fmt.Println("Custom armor removed:", name)
	} else {
		fmt.Println("Armor not found or is not custom. No action taken.")
	}
}

// verifyArmorType checks if the provided ArmorType is valid.
func verifyArmorType(t ArmorType) bool{
	switch t{
		case ArmorTypeLight, ArmorTypeMedium, ArmorTypeHeavy:
			return true
		default:
			return false
	}
}


func (a Armor) String() string {
	return fmt.Sprintf("%s (Type: %d, Defense: %d, Weight: %d, Rarity: %d)", a.Name, a.Type, a.Defense, a.Weight, a.Rarity)
}

func (a Armor) GetName() string {
	return a.Name
}

func (a Armor) GetDescription() string {
	return a.Description
}

func (a Armor) GetWeight() int {
	return a.Weight
}

func (a Armor) GetRarity() Rarity {
	return a.Rarity
}

func (a Armor) IsCustomEquipment() bool {
	return a.IsCustom
}
