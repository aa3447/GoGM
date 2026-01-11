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
	MaxHealth int
	Mana int
	MaxMana int
	Attributes PlayerAttributes
	Buffs map[string]gameLogic.Buff
	Inventory []equipment.Equipment
	CurrentArmor equipment.Armor
	CurrentWeapon equipment.Weapon
	PlayerPositionX int
	PlayerPositionY int
}

type PlayerAttributes struct{
	Strength int
	Dexterity int
	Intelligence int
	Constitution int
	Charisma int
	Wisdom int
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
	if p.Health > p.MaxHealth{
		p.Health = p.MaxHealth
	}
}

func (p *Player) RestoreMana(amount int){
	p.Mana += amount
	if p.Mana > p.MaxMana{
		p.Mana = p.MaxMana
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
	
}

func (p *Player) RemoveBuff(attr string){
	delete(p.Buffs, attr)
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
	fmt.Printf("Health: %d/%d, Mana: %d/%d\n", p.Health, p.MaxHealth, p.Mana, p.MaxMana)
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
