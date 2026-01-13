package gameLogic

import (
	"math/rand"
)

func DiceRoll(sides int) int{
	return rand.Intn(sides) + 1
}

func DiceRolls(sides int, rolls int) int{
	total := 0
	for range rolls{
		total += DiceRoll(sides)
	}
	return total
}

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

func DiceRollWithModifier(sides int, modifier int) int{
	return DiceRoll(sides) + modifier
}

func DiceRollsWithModifier(sides int, numberOfRolls int, modifier int) int{
	return DiceRolls(sides, numberOfRolls) + modifier
}

	


