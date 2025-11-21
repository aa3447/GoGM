package player

import (
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/equipment"
)

type Player struct{
	Name string
	Description string
	Background string
	Level int
	Experience int
	Health int
	MaxHealth int
	CurrentArmor equipment.Armor
	CurrentWeapon equipment.Weapon
	PlayerPositionX int
	PlayerPositionY int
}


func NewPlayer(name string, description string, background string) *Player{
	return &Player{
		Name: name,
		Description: description,
		Background: background,
		Level: 1,
		Experience: 0,
		Health: 100,
		MaxHealth: 100,
		PlayerPositionX: 0,
		PlayerPositionY: 0,
	}
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