package playerLogic

import (
	"fmt"
	"math/rand"
	"sort"

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
	Buffs map[string]Buff
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


// NewPlayer creates a new Player with the specified name, description, background, and stat generation method.
func NewPlayer(name, description, background, statMethod string, args ...[]int) *Player{
	var stats []int

	if statMethod == ""{
		statMethod = "roll"
	}
	switch statMethod {
		case "roll":
			stats = statMethodRoll()
		case "buy":
			stats = statMethodRoll() //placeholder for buy method
		case "assign":
			stats = args[0]
		default:
			fmt.Println("Invalid stat method, defaulting to roll")
			stats = statMethodRoll()
	}
	

	player := &Player{
		Name: name,
		Description: description,
		Background: background,
		Level: 1,
		Experience: 0,
		Health: 100,
		MaxHealth: 100,
		Attributes: PlayerAttributes{},
		PlayerPositionX: 0,
		PlayerPositionY: 0,
	}

	statAssign(player, stats)

	return player
}

// statMethodRoll generates stats by rolling 4d6 and dropping the lowest die for each of the 6 attributes.
func statMethodRoll() []int{
	stats := make([]int, 6)
	for i := range 6{
		rolls := make([]int, 4)
		for range 4{
			rolls = append(rolls, rand.Intn(6)+1)
		}
		//remove lowest roll
		sort.Ints(rolls)
		rollsLowestRemoved := rolls[1:]
		sum := 0
		for _, roll := range rollsLowestRemoved {
			sum += roll
		}
		stats[i] = sum
	}
	return stats
}

// statAssign assigns the generated stats to the player's attributes based on user input.
func statAssign(p *Player, stats []int){
	var choice int
	var choicesMade = make(map[int]bool)
	fmt.Printf("Here is your stat spread %v", stats)

	for i, stat := range stats{
		fmt.Printf("Stat %d: %d\n", i+1, stat)
		fmt.Print("Assign to (1) Strength, (2) Dexterity, (3) Intelligence, (4) Constitution, (5) Charisma, (6) Wisdom: ")
		
		
		fmt.Scan(&choice)
		choice = choiceValidation(choice, 1, 6)
		for choicesMade[choice]{
			fmt.Println("Stat already assigned, please choose another stat.")
			fmt.Scan(&choice)
			choice = choiceValidation(choice, 1, 6)
		}

		switch choice{
			case 1:
				p.Attributes.Strength = stat
			case 2:
				p.Attributes.Dexterity = stat
			case 3:
				p.Attributes.Intelligence = stat
			case 4:
				p.Attributes.Constitution = stat
			case 5:
				p.Attributes.Charisma = stat
			case 6:
				p.Attributes.Wisdom = stat
		}
		choicesMade[choice] = true
	}

	fmt.Printf("Final Attributes: %+v\n", p.Attributes)
}

func choiceValidation(input int, min int, max int) int{
	for input < min || input > max{
		fmt.Printf("Invalid choice. Please enter a number between %d and %d: ", min, max)
		_, err := fmt.Scan(&input)
		if err != nil{
			continue
		}
	}
	return input
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

func (p *PLayer) Heal(amount int){
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
	 buff := Buff{
		Attribute: attr,
		Amount: amount,
		Duration: duration,
	}
	
	if p.Buffs == nil{
		p.Buffs = make(map[string]Buff)
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

func (p *Player) GetBuffs() []Buff{
	buffs := []Buff{}
	for _, buff := range p.Buffs{
		buffs = append(buffs, buff)
	}
	return buffs
}

func (p *Player) GetBuffByAttribute(attr string) *Buff{
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
