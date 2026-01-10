package equipment

type Equipment interface{
	GetName() string
	GetDescription() string
	GetWeight() int
	GetRarity() Rarity
	IsCustomEquipment() bool
}