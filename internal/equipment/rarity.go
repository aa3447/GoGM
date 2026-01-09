package equipment

type Rarity int
const (
	Common Rarity  = iota
	Uncommon
	Rare
	Epic
	Legendary
)

// verifyRarity checks if the provided Rarity is valid.
func verifyRarity(r Rarity) bool{
	switch r{
		case Common, Uncommon, Rare, Epic, Legendary:
			return true
		default:
			return false
	}
}