package gameLogic

import (
	"slices"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/gameLogic"
)


func CombatLoop(players []*playerLogic.Player, enemies []*playerLogic.NPC){
	combatOrder := determineCombatOrder(players, enemies)

	for len(players) > 0 && len(enemies) > 0{
		for _, combatant := range combatOrder{
			if combatant[1] == 0{ // Player's turn
				player := players[combatant[2]]
				// Player action logic here
				_ = player // Placeholder to avoid unused variable error
			} else { // Enemy's turn
				enemy := enemies[combatant[2]]
				// Enemy action logic here
				_ = enemy // Placeholder to avoid unused variable error
			}
			// Check for end of combat conditions
			if len(players) == 0 || len(enemies) == 0{
				break
			}
		}
	}
}

func determineCombatOrder(players []*playerLogic.Player, enemies []*playerLogic.NPC) [][]int{
	// Store initiatives as [initiativeRoll, type(0=player,1=enemy), indexInOriginalSlice]
	var combatants [][]int
	var playerInitiatives [][]int
	var enemyInitiatives [][]int

	// Roll initiative for each combatant
	for i, c := range players{
		initRoll := gameLogic.DiceRollWithModifier(20, c.AttributeModifiers["Dexterity"])
		playerInitiatives = append(playerInitiatives, []int{initRoll, 0, i})
	}
	for i, c := range enemies{
		initRoll := gameLogic.DiceRollWithModifier(20, c.AttributeModifiers["Dexterity"])
		enemyInitiatives = append(enemyInitiatives, []int{initRoll, 1, i})
	}

	combatants = append(playerInitiatives, enemyInitiatives...)

	slices.SortFunc(combatants, func(a, b []int) int {
		return b[0] - a[0] // Descending order
	})

	return combatants
}