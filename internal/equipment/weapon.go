package equipment

import (
	"fmt"
)

type WeaponType int
const (
	MeleeWeaponTypeSword WeaponType  = iota
	MeleeWeaponTypeAxe
	MeleeWeaponTypeMace
	MeleeWeaponTypeDagger
	RangeWeaponTypeBow
	RangeWeaponTypeCrossbow
	RangeWeaponTypeSling
)



type Weapon struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Damage int `json:"damage"`
	Type WeaponType `json:"type"`
	Weight int `json:"weight"`
	Rarity Rarity `json:"rarity"`
	IsCustom bool `json:"is_custom"`
}

// Predefined weapon types
var mapOfWeapons = map[string]Weapon{
	"Short Sword": {
		Name: "Short Sword",
		Description: "A light and versatile melee weapon.",
		Damage: 6,
		Type: MeleeWeaponTypeSword,
		Weight: 5,
		Rarity: Common,
		IsCustom: false,
	},
	"Long Sword": {
		Name: "Long Sword",
		Description: "A balanced melee weapon with good reach.",
		Damage: 8,
		Type: MeleeWeaponTypeSword,
		Weight: 10,
		Rarity: Uncommon,
		IsCustom: false,
	},
	"Great Axe": {
		Name: "Great Axe",
		Description: "A heavy weapon that deals massive damage.",
		Damage: 12,
		Type: MeleeWeaponTypeAxe,
		Weight: 20,
		Rarity: Rare,
		IsCustom: false,
	},
	"Long Bow": {
		Name: "Long Bow",
		Description: "A ranged weapon effective at long distances.",
		Damage: 10,
		Type: RangeWeaponTypeBow,
		Weight: 15,
		Rarity: Uncommon,
		IsCustom: false,
	},
	"Short Bow": {
		Name: "Short Bow",
		Description: "A ranged weapon effective at short distances.",
		Damage: 8,
		Type: RangeWeaponTypeBow,
		Weight: 10,
		Rarity: Common,
		IsCustom: false,
	},
	"Crossbow": {
		Name: "Crossbow",
		Description: "A ranged weapon with high accuracy and power.",
		Damage: 12,
		Type: RangeWeaponTypeCrossbow,
		Weight: 20,
		Rarity: Rare,
		IsCustom: false,
	},
	"Sling": {
		Name: "Sling",
		Description: "A simple ranged weapon for throwing small projectiles.",
		Damage: 4,
		Type: RangeWeaponTypeSling,
		Weight: 3,
		Rarity: Common,
		IsCustom: false,
	},
}

// GetWeaponByName retrieves a Weapon by its name from the predefined weapon types.
func GetWeaponByName(name string) (Weapon, bool){
	weapon, exists := mapOfWeapons[name]
	return weapon, exists
}

// GetAllWeaponTypes returns a slice of all predefined weapon types.
func GetAllWeaponTypes() []Weapon{
	weapons := []Weapon{}
	for _, weapon := range mapOfWeapons{
		weapons = append(weapons, weapon)
	}
	return weapons
}

// GetAllWeapons returns a slice of all weapons, including custom ones.
func GetAllWeapons() []Weapon {
	weapons := []Weapon{}
	for _, weapon := range mapOfWeapons{
		weapons = append(weapons, weapon)
	}
	return weapons
}

// AddCustomWeapon adds a new custom weapon to the weapon map.
func AddCustomWeapon(name string, description string, damage int, t WeaponType, weight int, rarity Rarity){
	if !verifyWeaponType(t){
		return
	}
	if !verifyRarity(rarity){
		return
	}

	mapOfWeapons[name] = Weapon{
		Name: name,
		Description: description,
		Damage: damage,
		Type: t,
		Weight: weight,
		Rarity: rarity,
		IsCustom: true,
	}
}

// RemoveCustomWeapon removes a custom weapon from the weapon map.
func removeCustomWeapon(name string){
	weapon, exists := mapOfWeapons[name]
	if exists && weapon.IsCustom{
		delete(mapOfWeapons, name)
	} else {
		fmt.Println("Weapon not found or is not custom. No action taken.")
	}
}

// verifyWeaponType checks if the provided WeaponType is valid.
func verifyWeaponType(t WeaponType) bool{
	switch t{
		case MeleeWeaponTypeSword, MeleeWeaponTypeAxe, MeleeWeaponTypeMace, MeleeWeaponTypeDagger, RangeWeaponTypeBow, RangeWeaponTypeCrossbow, RangeWeaponTypeSling:
			return true
		default:
			return false
	}
}

func (w *Weapon) String() string {
	return fmt.Sprintf("%s (Type: %d, Damage: %d, Weight: %d, Rarity: %d)", w.Name, w.Type, w.Damage, w.Weight, w.Rarity)
}

func (w Weapon) GetName() string{
	return w.Name
}

func (w Weapon) GetDescription() string{
	return w.Description
}
func (w Weapon) GetWeight() int{
	return w.Weight
}
func (w Weapon) GetRarity() Rarity{
	return w.Rarity
}
func (w Weapon) IsCustomEquipment() bool{
	return w.IsCustom
}
