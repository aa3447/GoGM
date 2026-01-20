package playerLogic

import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/classes"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/equipment"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/gameLogic"
)

type Player struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Background string `json:"background"`
	Class string `json:"class"`
	Level int `json:"level"`
	LevelTrack LevelingTrack `json:"level_track"`
	Experience int `json:"experience"`
	Health int `json:"health"`
	Mana int `json:"mana"`
	Attributes PlayerAttributes `json:"attributes"`
	Buffs map[string]gameLogic.Buff `json:"buffs"`
	Inventory []equipment.Equipment `json:"inventory"`
	AttributeModifiers map[string]int `json:"attributeModifiers"`
	DerivedStats map[string]map[string]int `json:"derivedStats"`
	HitDieRolls []int `json:"hitDieRolls"`
	HitDieTotal int `json:"hitDieTotal"`
	CurrentArmor equipment.Armor `json:"currentArmor"`
	CurrentWeapon equipment.Weapon `json:"currentWeapon"`
	PlayerPositionX int `json:"playerPositionX"`
	PlayerPositionY int `json:"playerPositionY"`
	Initiative int `json:"initiative"`
	IsUnconscious bool `json:"isUnconscious"`
	IsAlive bool `json:"isAlive"`
}

type PlayerAttributes struct{
	Strength int `json:"strength"`
	Dexterity int `json:"dexterity"`
	Intelligence int `json:"intelligence"`
	Constitution int `json:"constitution"`
	Charisma int `json:"charisma"`
	Wisdom int `json:"wisdom"`
}

// SetAttributeModifiers calculates and sets the attribute modifiers based on the player's attributes.
func (p *Player) SetAttributeModifiers(){
	p.AttributeModifiers = map[string]int{
		"Strength": (p.Attributes.Strength - 10) / 2,
		"Dexterity": (p.Attributes.Dexterity - 10) / 2,
		"Intelligence": (p.Attributes.Intelligence - 10) / 2,
		"Constitution": (p.Attributes.Constitution - 10) / 2,
		"Charisma": (p.Attributes.Charisma - 10) / 2,
		"Wisdom": (p.Attributes.Wisdom - 10) / 2,
	}
}

// SetDerivedStats calculates and sets the derived stats based on the player's attribute modifiers.
func (p *Player) SetDerivedStats(){
	p.DerivedStats = map[string]map[string]int{
		"Strength": {
			"MeleeToHit": p.AttributeModifiers["Strength"] * p.Level,
			"MeleeAttackPower": p.AttributeModifiers["Strength"] * p.Level * 2,
		},
		"Dexterity": {
			"RangedToHit": p.AttributeModifiers["Dexterity"] * p.Level,
			"RangedAttackPower": p.AttributeModifiers["Dexterity"] * p.Level * 2,
			"BaseArmor": p.AttributeModifiers["Dexterity"] * p.Level,
		},
		"Intelligence": {
			"ShortRestManaRegen": p.AttributeModifiers["Intelligence"] * p.Level,
			"MaxMana": (p.AttributeModifiers["Intelligence"] * p.Level * 3) + classes.GetClassManaBonus(p.Class, p.Level),
			"SpellToHit": p.AttributeModifiers["Intelligence"] * p.Level,
			"SpellPower": p.AttributeModifiers["Intelligence"] * p.Level * 3,
		},
		"Constitution": {
			"ShortRestHealthRegen": p.AttributeModifiers["Constitution"] * 2 * p.Level,
			"MaxHealth": (p.AttributeModifiers["Constitution"] * p.Level * 5) + p.HitDieTotal,
		},
	}
}

// SetDerivedStatsForAttribute updates the derived stats for a specific attribute.
func (p *Player) SetDerivedStatsForAttribute(attribute string){
	switch attribute{
		case "Strength":
			p.DerivedStats[attribute] = map[string]int{
				"MeleeToHit": p.AttributeModifiers[attribute] * p.Level,
				"MeleeAttackPower": p.AttributeModifiers[attribute] * p.Level * 2,
			}
		case "Dexterity":
			p.DerivedStats[attribute] = map[string]int{
				"RangedToHit": p.AttributeModifiers[attribute] * p.Level,
				"RangedAttackPower": p.AttributeModifiers[attribute] * p.Level * 2,
				"BaseArmor": p.AttributeModifiers[attribute] * p.Level,
			}
		case "Intelligence":
			p.DerivedStats[attribute] = map[string]int{
				"ShortRestManaRegen": p.AttributeModifiers[attribute] * p.Level,
				"MaxMana": (p.AttributeModifiers[attribute]  * p.Level * 3) + classes.GetClassManaBonus(p.Class, p.Level),
				"SpellToHit": p.AttributeModifiers[attribute] * p.Level,
				"SpellPower": p.AttributeModifiers[attribute] * p.Level * 3,
			}
		case "Constitution":
			p.DerivedStats[attribute] = map[string]int{
				"ShortRestHealthRegen": p.AttributeModifiers[attribute] * 2 * p.Level,
				"MaxHealth": (p.AttributeModifiers[attribute] * p.Level * 5) + p.HitDieTotal,
			}
	}
}

// SetDerivedStat sets a specific derived stat for the player.
func (p *Player) SetDerivedStat(attribute string, stat string, value int){
	if _, exists := p.DerivedStats[attribute]; exists{
		p.DerivedStats[attribute][stat] = value
	} else{
		p.DerivedStats[attribute] = map[string]int{
			stat: value,
		}
	}
}


// GetDerivedStat returns the value of a derived stat.
func (p *Player) GetDerivedStat(attribute string, stat string) int{
	if value, exists := p.DerivedStats[attribute]; exists{
		return value[stat]
	}
	return 0
}

func (p *Player) RollHitDie(){
	hitDieType := classes.Classes[p.Class].HitDie
	fmt.Println("Use max roll for hit die? (y/n): ")
	var input string
	fmt.Scan(&input)
	
	var roll int
	if input == "y" || input == "Y"{
		roll = hitDieType
	} else {
		roll = gameLogic.DiceRoll(hitDieType)
	}
	p.HitDieRolls = append(p.HitDieRolls, roll)
	p.HitDieTotal += roll
}
	

// SetLocation sets the player's position on the map.
func (p *Player) SetLocation(y int, x int){
	p.PlayerPositionY = y
	p.PlayerPositionX = x
}

// Move updates the player's position by the specified deltas.
func (p *Player) Move(deltaY int, deltaX int){
	p.PlayerPositionY += deltaY
	p.PlayerPositionX += deltaX
}

// GetPlayerPosition returns the player's current position on the map.
func (p *Player) GetPlayerPosition() (int, int){
	return p.PlayerPositionY, p.PlayerPositionX
}

// Heal increases the player's health by the specified amount, up to the maximum health.
func (p *Player) Heal(amount int){
	p.Health += amount
	if p.Health > p.DerivedStats["Constitution"]["MaxHealth"]{
		p.Health = p.DerivedStats["Constitution"]["MaxHealth"]
	}
}

// RestoreMana increases the player's mana by the specified amount, up to the maximum mana.
func (p *Player) RestoreMana(amount int){
	p.Mana += amount
	if p.Mana > p.DerivedStats["Intelligence"]["MaxMana"]{
		p.Mana = p.DerivedStats["Intelligence"]["MaxMana"]
	}
}

// TakeDamage decreases the player's health by the specified amount.
func (p *Player) TakeDamage(amount int){
	p.Health -= amount
	if p.Health < 0{
		p.Health = 0
	}
}

// UseMana decreases the player's mana by the specified amount if enough mana is available.
func (p *Player) UseMana(amount int) bool{
	if p.Mana < amount{
		return false
	}
	p.Mana -= amount
	return true
}

// BuffAttribute applies a buff to the specified attribute for a given duration.
func (p *Player) BuffAttribute(attr string, amount int, duration int){
	 buff := gameLogic.Buff{
		Attribute: attr,
		Amount: amount,
		Duration: duration,
	}
	
	if p.Buffs == nil{
		p.Buffs = make(map[string]gameLogic.Buff)
	}
	if existingBuff, exists := p.Buffs[attr]; exists{
		if amount > existingBuff.Amount{
			existingBuff.Duration = existingBuff.Duration + duration/2
			existingBuff.Amount = amount
			p.Buffs[attr] = existingBuff
		} else if amount == existingBuff.Amount && duration > existingBuff.Duration{
			existingBuff.Duration += duration
			p.Buffs[attr] = existingBuff
		}
		return
	}
	p.Buffs[attr] = buff
	p.SetDerivedStatsForAttribute(attr)
	
}

// RemoveBuff removes the buff from the specified attribute.
func (p *Player) RemoveBuff(attr string){
	delete(p.Buffs, attr)
	p.SetDerivedStatsForAttribute(attr)
}

// GetBuffs returns a slice of all active buffs on the player.
func (p *Player) GetBuffs() []gameLogic.Buff{
	buffs := []gameLogic.Buff{}
	for _, buff := range p.Buffs{
		buffs = append(buffs, buff)
	}
	return buffs
}

// GetBuffByAttribute returns the buff for a specific attribute, if it exists.
func (p *Player) GetBuffByAttribute(attr string) *gameLogic.Buff{
	if buff, exists := p.Buffs[attr]; exists{
		return &buff
	}
	return nil
}

// ListBuffs prints all active buffs on the player.
func (p *Player) ListBuffs(){
	for _, buff := range p.Buffs{
		fmt.Println(buff.String())
	}
}

// ShowInventory prints the player's inventory.
func (p *Player) ShowInventory(){
	fmt.Println("Inventory:")
	for _, item := range p.Inventory{
		fmt.Printf("- %s: %s\n", item.GetName(), item.GetDescription())
	}
}

// ShowStats prints the player's stats.
func (p *Player) ShowStats(){
	fmt.Printf("Player: %s\n", p.Name)
	fmt.Printf("Level: %d, Experience: %d\n", p.Level, p.Experience)
	fmt.Printf("Health: %d/%d, Mana: %d/%d\n", p.Health, p.DerivedStats["MaxHealth"], p.Mana, p.DerivedStats["MaxMana"])
	fmt.Printf("Attributes:\n")
	fmt.Printf("  Strength: %d\n", p.Attributes.Strength)
	fmt.Printf("  Dexterity: %d\n", p.Attributes.Dexterity)
	fmt.Printf("  Intelligence: %d\n", p.Attributes.Intelligence)
	fmt.Printf("  Constitution: %d\n", p.Attributes.Constitution)
	fmt.Printf("  Charisma: %d\n", p.Attributes.Charisma)
	fmt.Printf("  Wisdom: %d\n", p.Attributes.Wisdom)
}


func (p *Player) AddEquipmentToInventory(item equipment.Equipment){
	p.Inventory = append(p.Inventory, item)
}

func (p *Player) EquipWeapon(weapon equipment.Weapon){
	p.CurrentWeapon = weapon
}

func (p *Player) EquipArmor(armor equipment.Armor){
	p.CurrentArmor = armor
}

// ShowEquipment prints the player's currently equipped weapon and armor.
func (p *Player) ShowEquipment(){
	fmt.Println("Current Equipment:")
	fmt.Printf("  Weapon: %s\n", p.CurrentWeapon.String())
	fmt.Printf("  Armor: %s\n", p.CurrentArmor.String())
}
