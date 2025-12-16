package gameLogic

import (
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
)

type Gamestate struct{
	Players map[string]*playerLogic.Player `json:"players"`
}


