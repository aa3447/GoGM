package classes

type Class struct{
	Name string `json:"name"`
	Description string `json:"description"`
	HitDie int `json:"hitDie"`
	ManaBonusPerLevel int `json:"manaBonusPerLevel"`
	PrimaryAttributes []string `json:"primaryAttributes"`
	Proficiencies []string `json:"proficiencies"`
	StartingEquipment []string `json:"startingEquipment"`
	Skills []string `json:"skills"`
}

var Classes = map[string]Class{
	"Warrior": {
		Name: "Warrior",
		Description: "A strong and resilient fighter.",
		HitDie: 10,
		ManaBonusPerLevel: 0,
		PrimaryAttributes: []string{"Strength", "Constitution"},
		Proficiencies: []string{"All Armor", "Shields", "Melee Weapons"},
		StartingEquipment: []string{"Longsword", "Shield", "Chainmail"},
		Skills: []string{"Power Attack", "Shield Bash"},
	},
	"Mage": {
		Name: "Mage",
		Description: "A master of arcane arts.",
		HitDie: 6,
		ManaBonusPerLevel: 8,
		PrimaryAttributes: []string{"Intelligence", "Wisdom"},
		Proficiencies: []string{"Light Armor", "Staves", "Wands"},
		StartingEquipment: []string{"Spellbook", "Staff", "Robe"},
		Skills: []string{"Fireball", "Teleport"},
	},
	"Rogue": {
		Name: "Rogue",
		Description: "A stealthy and agile character.",
		HitDie: 8,
		ManaBonusPerLevel: 2,
		PrimaryAttributes: []string{"Dexterity", "Charisma"},
		Proficiencies: []string{"Light Armor", "Daggers", "Bows"},
		StartingEquipment: []string{"Dagger", "Bow", "Leather Armor"},
		Skills: []string{"Backstab", "Stealth"},
	},
	"Cleric": {
		Name: "Cleric",
		Description: "A holy warrior and healer.",
		HitDie: 8,
		ManaBonusPerLevel: 4,
		PrimaryAttributes: []string{"Wisdom", "Charisma"},
		Proficiencies: []string{"Medium Armor", "Shields", "Maces"},
		StartingEquipment: []string{"Mace", "Shield", "Chainmail"},
		Skills: []string{"Heal", "Turn Undead"},
	},
	"Ranger": {
		Name: "Ranger",
		Description: "A skilled hunter and tracker.",
		HitDie: 10,
		ManaBonusPerLevel: 2,
		PrimaryAttributes: []string{"Dexterity", "Wisdom"},
		Proficiencies: []string{"Medium Armor", "Bows", "Dual Wielding"},
		StartingEquipment: []string{"Longbow", "Dual Daggers", "Leather Armor"},
		Skills: []string{"Multiattack", "Tracking"},
	},
}	

var ClassList []string

func init(){
	for className := range Classes{
		ClassList = append(ClassList, className)
	}
}

// GetClassManaBonus returns the mana bonus for a given class and level.
func GetClassManaBonus(className string, level int) int {
	if class, exists := Classes[className]; exists {
		return class.ManaBonusPerLevel * level
	}
	return 0
}