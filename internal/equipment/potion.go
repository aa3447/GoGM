package equipment

type PotionEffect int
const (
	EffectHeal PotionEffect = iota
	EffectMana
	EffectStrength
	EffectDexterity
	EffectIntelligence
)

// Potion represents a consumable item that grants various effects.
type Potion struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Effect PotionEffect `json:"effect"`
	Power int `json:"power"` // Magnitude of the effect
	Rarity Rarity `json:"rarity"`
	Duration int `json:"duration"` // Duration in seconds for temporary effects; 0 for instant effects
	IsCustom bool `json:"is_custom"`
}

// Predefined potion types
var mapOfPotions = map[string]Potion{
	"Healing Potion": {
		Name: "Healing Potion",
		Description: "Restores a small amount of health.",
		Effect: EffectHeal,
		Power: 20,
		Rarity: Common,
		Duration: 0,
		IsCustom: false,
	},
	"Mana Potion": {
		Name: "Mana Potion",
		Description: "Restores a small amount of mana.",
		Effect: EffectMana,
		Power: 15,
		Rarity: Common,
		Duration: 0,
		IsCustom: false,
	},
	"Strength Potion": {
		Name: "Strength Potion",
		Description: "Temporarily increases strength.",
		Effect: EffectStrength,
		Power: 5,
		Rarity: Uncommon,
		Duration: 60,
		IsCustom: false,
	},
}

// GetPotionByName retrieves a Potion by its name from the predefined potion types.
func GetPotionByName(name string) (Potion, bool){
	potion, exists := mapOfPotions[name]
	return potion, exists
}

// GetAllPotionTypes returns a slice of all predefined potion types.
func GetAllPotionTypes() []Potion{
	potions := make([]Potion, 0, len(mapOfPotions))
	for _, potion := range mapOfPotions{
		potions = append(potions, potion)
	}
	return potions
}

func (p *Potion) GetName() string{
	return p.Name
}
func (p *Potion) GetDescription() string{
	return p.Description
}
func (p *Potion) GetWeight() int{
	return 1 // Potions have a fixed weight of 1 unit
}
func (p *Potion) GetRarity() Rarity{
	return p.Rarity
}
func (p *Potion) IsCustomEquipment() bool{
	return p.IsCustom
}
