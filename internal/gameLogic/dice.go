package gameLogic

import (
	"math/rand"
)

// DiceRoll simulates rolling a single die with the specified number of sides.
func DiceRoll(sides int) int{
	return rand.Intn(sides) + 1
}

// DiceRolls simulates rolling multiple dice with the specified number of sides and returns the total sum.
func DiceRolls(sides int, rolls int) int{
	total := 0
	for range rolls{
		total += DiceRoll(sides)
	}
	return total
}

// DiceRollMultipleDice rolls the specified number of dice with the given sides and returns a slice of individual rolls along with their sum as the last element.
func DiceRollMultipleDice(sides int, numberOfRolls int) []int{
	rolls := make([]int, numberOfRolls)
	sum := 0
	for i := range rolls{
		rolls[i] = DiceRoll(sides)
		sum += rolls[i]
	}
	rolls = append(rolls, sum)
	return rolls
}

// DiceRollWithModifier rolls a die with the specified number of sides and adds a modifier to the result.
func DiceRollWithModifier(sides int, modifier int) int{
	return DiceRoll(sides) + modifier
}

// DiceRollsWithModifier rolls multiple dice with the specified number of sides, sums the results, and adds a modifier.
func DiceRollsWithModifier(sides int, numberOfRolls int, modifier int) int{
	return DiceRolls(sides, numberOfRolls) + modifier
}

	


