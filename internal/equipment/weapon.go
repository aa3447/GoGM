package equipment

import (
	"fmt"
)

type WeaponType int
const (
	WeaponTypeSword WeaponType  = iota
	WeaponTypeAxe
	WeaponTypeMace
	WeaponTypeBow
	WeaponTypeDagger
)

type Weapon struct{
	Name string
	Description string
	Damage int
	Type WeaponType
	Weight int
	Rarity Rarity
	IsCustom bool
}

// Predefined weapon types
var mapOfWeapons = map[string]Weapon{
	"Short Sword": {
		Name: "Short Sword",
		Description: "A light and versatile melee weapon.",
		Damage: 6,
		Type: WeaponTypeSword,
		Weight: 5,
		Rarity: Common,
		IsCustom: false,
	},
	"Long Sword": {
		Name: "Long Sword",
		Description: "A balanced melee weapon with good reach.",
		Damage: 8,
		Type: WeaponTypeSword,
		Weight: 10,
		Rarity: Uncommon,
		IsCustom: false,
	},
	"Great Axe": {
		Name: "Great Axe",
		Description: "A heavy weapon that deals massive damage.",
		Damage: 12,
		Type: WeaponTypeAxe,
		Weight: 20,
		Rarity: Rare,
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
		case WeaponTypeSword, WeaponTypeAxe, WeaponTypeMace, WeaponTypeBow, WeaponTypeDagger:
			return true
		default:
			return false
	}
}

func (w Weapon) String() string {
	return fmt.Sprintf("%s (Type: %d, Damage: %d, Weight: %d, Rarity: %d)", w.Name, w.Type, w.Damage, w.Weight, w.Rarity)
}