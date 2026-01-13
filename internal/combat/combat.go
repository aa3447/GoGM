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
			switch v := combatant.(type){
				case *playerLogic.Player:
					// Player's turn logic here
				case *playerLogic.NPC:
					// NPC's turn logic here
			}
		}
	}
}

func determineCombatOrder(players []*playerLogic.Player, enemies []*playerLogic.NPC) []any{
	var combatants []any
	for _, p := range players{
		combatants = append(combatants, p)
	}
	for _, e := range enemies{
		combatants = append(combatants, e)
	}

	// Roll initiative for each combatant
	for _, c := range combatants{
		switch v := c.(type){
			case *playerLogic.Player:
				initiativeRoll := gameLogic.DiceRoll(20) + v.AttributeModifiers["Dexterity"]
				v.Initiative = initiativeRoll
			case *playerLogic.NPC:
				initiativeRoll := gameLogic.DiceRoll(20) + v.AttributeModifiers["Dexterity"]
				v.Initiative = initiativeRoll
		}
	}

	// Sort combatants by initiative
	slices.SortFunc(combatants, func(a, b any) int {
		var initA, initB int
		switch v := a.(type){
			case *playerLogic.Player:
				initA = v.Initiative
			case *playerLogic.NPC:
				initA = v.Initiative
		}
		switch v := b.(type){
			case *playerLogic.Player:
				initB = v.Initiative
			case *playerLogic.NPC:
				initB = v.Initiative
		}
		return initB - initA // Descending order
	})

	return combatants
}