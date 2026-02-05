package equipment

type Equipment interface{
	GetName() string
	GetDescription() string
	GetWeight() int
	GetRarity() Rarity
	IsCustomEquipment() bool
	String() string
}

// GetEquipmentByName retrieves equipment by its name.
func GetEquipmentByName(name string) (Equipment, bool){
	if weapon, exists := mapOfWeapons[name]; exists{
		return weapon, true
	}
	if armor, exists := mapOfArmors[name]; exists{
		return armor, true
	}
	return nil, false
}