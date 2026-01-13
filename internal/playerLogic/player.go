package playerLogic

import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/equipment"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/gameLogic"
)

type Player struct{
	Name string
	Description string
	Background string
	Class string
	Level int
	Experience int
	Health int
	Mana int
	Attributes PlayerAttributes
	Buffs map[string]gameLogic.Buff
	Inventory []equipment.Equipment
	AttributeModifiers map[string]int
	DerivedStats map[string]map[string]int
	CurrentArmor equipment.Armor
	CurrentWeapon equipment.Weapon
	PlayerPositionX int
	PlayerPositionY int
	Initiative int
	IsUnconscious bool
	IsAlive bool
}

type PlayerAttributes struct{
	Strength int
	Dexterity int
	Intelligence int
	Constitution int
	Charisma int
	Wisdom int
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
			"MeleeToHit": p.AttributeModifiers["Strength"],
			"MeleeAttackPower": p.AttributeModifiers["Strength"] * 2,
		},
		"Dexterity": {
			"RangedToHit": p.AttributeModifiers["Dexterity"],
			"RangedAttackPower": p.AttributeModifiers["Dexterity"] * 2,
			"BaseArmor": p.AttributeModifiers["Dexterity"],
		},
		"Intelligence": {
			"ShortRestManaRegen": p.AttributeModifiers["Intelligence"],
			"MaxMana": p.AttributeModifiers["Intelligence"] * 5,
			"SpellToHit": p.AttributeModifiers["Intelligence"],
			"SpellPower": p.AttributeModifiers["Intelligence"] * 3,
		},
		"Constitution": {
			"ShortRestHealthRegen": p.AttributeModifiers["Constitution"] * 2,
			"MaxHealth": p.AttributeModifiers["Constitution"] * 10,
		},
	}
}

// SetDerivedStatsForAttribute updates the derived stats for a specific attribute.
func (p *Player) SetDerivedStatsForAttribute(attribute string){
	switch attribute{
		case "Strength":
			p.DerivedStats[attribute] = map[string]int{
				"MeleeToHit": p.AttributeModifiers[attribute],
				"MeleeAttackPower": p.AttributeModifiers[attribute] * 2,
			}
		case "Dexterity":
			p.DerivedStats[attribute] = map[string]int{
				"RangedToHit": p.AttributeModifiers[attribute],
				"RangedAttackPower": p.AttributeModifiers[attribute] * 2,
				"BaseArmor": p.AttributeModifiers[attribute],
			}
		case "Intelligence":
			p.DerivedStats[attribute] = map[string]int{
				"ShortRestManaRegen": p.AttributeModifiers[attribute],
				"MaxMana": p.AttributeModifiers[attribute] * 5,
				"SpellToHit": p.AttributeModifiers[attribute],
				"SpellPower": p.AttributeModifiers[attribute] * 3,
			}
		case "Constitution":
			p.DerivedStats[attribute] = map[string]int{
				"ShortRestHealthRegen": p.AttributeModifiers[attribute] * 2,
				"MaxHealth": p.AttributeModifiers[attribute] * 10,
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


func (p *Player) SetLocation(y int, x int){
	p.PlayerPositionY = y
	p.PlayerPositionX = x
}

func (p *Player) Move(deltaY int, deltaX int){
	p.PlayerPositionY += deltaY
	p.PlayerPositionX += deltaX
}

func (p *Player) GetPlayerPosition() (int, int){
	return p.PlayerPositionY, p.PlayerPositionX
}

func (p *Player) Heal(amount int){
	p.Health += amount
	if p.Health > p.DerivedStats["Constitution"]["MaxHealth"]{
		p.Health = p.DerivedStats["Constitution"]["MaxHealth"]
	}
}

func (p *Player) RestoreMana(amount int){
	p.Mana += amount
	if p.Mana > p.DerivedStats["Intelligence"]["MaxMana"]{
		p.Mana = p.DerivedStats["Intelligence"]["MaxMana"]
	}
}

func (p *Player) TakeDamage(amount int){
	p.Health -= amount
	if p.Health < 0{
		p.Health = 0
	}
}

func (p *Player) UseMana(amount int) bool{
	if p.Mana < amount{
		return false
	}
	p.Mana -= amount
	return true
}

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

func (p *Player) RemoveBuff(attr string){
	delete(p.Buffs, attr)
	p.SetDerivedStatsForAttribute(attr)
}

func (p *Player) GetBuffs() []gameLogic.Buff{
	buffs := []gameLogic.Buff{}
	for _, buff := range p.Buffs{
		buffs = append(buffs, buff)
	}
	return buffs
}

func (p *Player) GetBuffByAttribute(attr string) *gameLogic.Buff{
	if buff, exists := p.Buffs[attr]; exists{
		return &buff
	}
	return nil
}

func (p *Player) ListBuffs(){
	for _, buff := range p.Buffs{
		fmt.Println(buff.String())
	}
}

func (p *Player) ShowInventory(){
	fmt.Println("Inventory:")
	for _, item := range p.Inventory{
		fmt.Printf("- %s: %s\n", item.GetName(), item.GetDescription())
	}
}

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

func (p *Player) ShowEquipment(){
	fmt.Println("Current Equipment:")
	fmt.Printf("  Weapon: %s\n", p.CurrentWeapon.String())
	fmt.Printf("  Armor: %s\n", p.CurrentArmor.String())
}
