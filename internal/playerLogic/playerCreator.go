package playerLogic

import (
	"fmt"
	"math/rand"
	"sort"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/equipment"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/gameLogic"
)

// CreatePlayer interacts with the user to create a new Player character.
func CreatePlayer() *Player{
	fmt.Println("Creating new player.")
	fmt.Print("Enter player name: ")
	name := io.GetInput()[0]
	fmt.Print("Enter player description: ")
	description := io.GetInput()[0]
	fmt.Print("Enter player background: ")
	background := io.GetInput()[0]
	fmt.Print("Enter stat generation method (roll, buy, assign): ")
	statMethod := io.GetInput()[0]
	var stats []int
	
	if statMethod == "assign"{
		var acceptedStats = false
		for !acceptedStats{
			fmt.Print("Enter stats as space-separated integers (Str Dex Int Con Cha Wis): ")
			statInputs := io.GetInput()
			for _, s := range statInputs{
				var stat int
				fmt.Sscanf(s, "%d", &stat)
				stats = append(stats, stat)
			}
			
			fmt.Printf("You entered the following stats: %v\n", stats)
			fmt.Print("Are these correct? (y/n): ")
			confirmation := io.GetInput()[0]
			if confirmation == "y" || confirmation == "Y"{
				acceptedStats = true
			} else{
				stats = []int{}
			}
		}
	}
	
	player := NewPlayer(name, description, background, statMethod, stats)
	fmt.Println("Player created:", player.Name)
	return player
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
			fmt.Printf("Rolled stats: %v\n", stats)
			fmt.Println("Do you want to reroll? (y/n): ")
			
			var rerollChoice string
			var reStats []int
			fmt.Scan(&rerollChoice)
			if rerollChoice == "y" || rerollChoice == "Y"{
				reStats := statMethodRoll()
				fmt.Printf("Rolled stats: %v\n", reStats)
			}

			if len(reStats) > 0{
				fmt.Println("Do you want to use your original or re-rolled stats? (0 for original, 1 for new)")
				var finalChoice int
				fmt.Scan(&finalChoice)
				finalChoice = choiceValidation(finalChoice, 0, 1)
				if finalChoice == 1{
					stats = reStats
				}
			}
		case "buy":
			stats = pointBuySystem()
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
		Class: "",
		Level: 1,
		Experience: 0,
		Health: 100,
		Mana: 50,
		Attributes: PlayerAttributes{},
		Buffs: make(map[string]gameLogic.Buff),
		Inventory: []equipment.Equipment{},
		DerivedStats: make(map[string]map[string]int),
		CurrentArmor: equipment.Armor{},
		CurrentWeapon: equipment.Weapon{},
		PlayerPositionX: 0,
		PlayerPositionY: 0,
	}

	statAssign(player, stats)

	player.SetDerivedStats()
	player.Health = player.DerivedStats["Constitution"]["MaxHealth"]
	player.Mana = player.DerivedStats["Intelligence"]["MaxMana"]

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

// pointBuySystem allows the user to allocate points to stats using a point-buy system.
func pointBuySystem() []int{
	var acceptedStats = false
	points := 20
	stats := [6]int{10,10,10,10,10,10}
	cost := map[int]int{
		7: -4,
		8: -2,
		9: -1,
		10: 0,
		11: 1,
		12: 2,
		13: 3,
		14: 5,
		15: 7,
		16: 10,
		17: 13,
		18: 17,
	}

	fmt.Printf("Points: %d\n", points)
	fmt.Printf("Starting stats (Str Dex Int Con Cha Wis): %v\n", stats)

	for acceptedStats == false{
		fmt.Print("Enter command (costs, stats, set, done): ")
		var command string
		fmt.Scan(&command)
		switch command{
			case "costs":
				fmt.Printf("Costs: %v\n", cost)
			case "stats":
				fmt.Printf("Current stats (Str Dex Int Con Cha Wis): %v\n", stats)
			case "set":
				var attrIndex int
				var newValue int
				fmt.Print("Enter attribute index (1-6) to set: ")
				fmt.Scan(&attrIndex)
				attrIndex = choiceValidation(attrIndex, 1, 6) - 1
				
				fmt.Print("Enter new value (7-18): ")
				fmt.Scan(&newValue)
				newValue = choiceValidation(newValue, 7, 18)
				
				currentValue := stats[attrIndex]
				pointChange := cost[newValue] - cost[currentValue]
				if points - pointChange < 0{
					fmt.Println("Not enough points for this change.")
				} else{
					stats[attrIndex] = newValue
					points -= pointChange
					fmt.Printf("Points remaining: %d\n", points)
				}
			case "done":
				if points >= 0{
					fmt.Printf("You still have %d points remaining. Are you sure you want to finish? (y/n): ", points)
					var finishChoice string
					fmt.Scan(&finishChoice)
					if finishChoice == "y" || finishChoice == "Y"{
						acceptedStats = true
					}
				} else {
					fmt.Printf("You have overspent by %d points. Please adjust your stats.\n", -points)
				}
			default:
				fmt.Println("Unknown command. Available commands: costs, stats, set, done")
		}
	}
	return stats[:]
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

// choiceValidation ensures the input is within the specified range.
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

