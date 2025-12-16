package playerLogic

import (
	"fmt"
	"math/rand"
	"sort"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/equipment"
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
	Attributes PlayerAttributes
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