package gameLogic

import (
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
)

// Gamestate represents the current state of the game, including players and other dynamic elements.
type Gamestate struct{
	Players map[string]*playerLogic.Player `json:"players"`
}


