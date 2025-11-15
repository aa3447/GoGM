package main

import (
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
)

func main(){
	gameMap, err := mapLogic.GenRandomMap(5, 5, 0.2, []float64{0.5, 0.2, 0.2, 0.1})
	if err != nil{
		panic(err)
	}
	mapLogic.PrintMap(gameMap)
	mapLogic.PrintMapDebug(gameMap)
}